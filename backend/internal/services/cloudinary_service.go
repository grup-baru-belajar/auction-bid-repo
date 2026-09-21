package services

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/config"
)

const (
	MaxImageBytes  = 1 << 20 // 1MB
	maxImageWidth  = 4000
	maxImageHeight = 4000
	uploadFolder   = "auction-bid"
	maxRespBytes   = 1 << 20
)

var publicIDPattern = regexp.MustCompile(`^` + uploadFolder + `/[a-f0-9]{32}$`)

// ---------------------------------------------------------------------------
// Error types
// ---------------------------------------------------------------------------

type UploadError struct {
	Code    string
	Message string
	Err     error
}

func (e *UploadError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *UploadError) Unwrap() error { return e.Err }

func (e *UploadError) Is(target error) bool {
	t, ok := target.(*UploadError)
	return ok && t.Code == e.Code
}

func (e *UploadError) WithCause(err error) *UploadError {
	return &UploadError{Code: e.Code, Message: e.Message, Err: err}
}

func newUploadError(code, msg string) *UploadError {
	return &UploadError{Code: code, Message: msg}
}

var (
	ErrInvalidRequest    = newUploadError("INVALID_REQUEST", "request must be multipart/form-data")
	ErrImageRequired     = newUploadError("IMAGE_REQUIRED", "image file is required (form field: image)")
	ErrEmptyFile         = newUploadError("EMPTY_FILE", "image file is empty")
	ErrImageTooLarge     = newUploadError("FILE_TOO_LARGE", "image file too large, max 1MB")
	ErrUnsupportedFormat = newUploadError("UNSUPPORTED_FORMAT", "unsupported image format, only PNG, JPEG, WEBP allowed")
	ErrInvalidImageFile  = newUploadError("INVALID_IMAGE", "file is corrupted or not a valid image")
	ErrImageDimensions   = newUploadError("IMAGE_DIMENSIONS", "image dimensions too large, max 4000x4000")
	ErrInvalidPublicID   = newUploadError("INVALID_PUBLIC_ID", "public_id is invalid")

	ErrUploadBlocked  = newUploadError("NETWORK_BLOCKED", "connection to image storage was blocked by a network proxy/firewall (e.g. Zscaler)")
	ErrUploadTimeout  = newUploadError("UPLOAD_TIMEOUT", "image storage did not respond in time")
	ErrUploadNetwork  = newUploadError("NETWORK_ERROR", "cannot reach image storage")
	ErrStorageAuth    = newUploadError("STORAGE_CONFIG_ERROR", "image storage is not configured correctly")
	ErrStorageLimited = newUploadError("STORAGE_RATE_LIMITED", "image storage rate limit reached, try again later")
	ErrStorageReject  = newUploadError("STORAGE_REJECTED", "image storage rejected the file")
	ErrUploadFailed   = newUploadError("UPLOAD_FAILED", "failed to upload image to cloud storage")
	ErrDeleteFailed   = newUploadError("DELETE_FAILED", "failed to delete image from cloud storage")
	ErrImageNotFound  = newUploadError("IMAGE_NOT_FOUND", "image not found in cloud storage")
)

// ---------------------------------------------------------------------------
// Service
// ---------------------------------------------------------------------------

type UploadResult struct {
	URL      string
	PublicID string
	Format   string
	Width    int
	Height   int
	Bytes    int
}

type CloudinaryService interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (*UploadResult, error)
	Delete(ctx context.Context, publicID string) error
}

type cloudinaryService struct {
	cloudName string
	apiKey    string
	apiSecret string
	client    *http.Client
}

func NewCloudinaryService(cfg config.CloudinaryConfig) CloudinaryService {
	return &cloudinaryService{
		cloudName: cfg.CloudName,
		apiKey:    cfg.APIKey,
		apiSecret: cfg.APISecret,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

type cloudinaryUploadResponse struct {
	SecureURL string `json:"secure_url"`
	PublicID  string `json:"public_id"`
	Format    string `json:"format"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	Bytes     int    `json:"bytes"`
	Error     *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type cloudinaryDeleteResponse struct {
	Result string `json:"result"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (s *cloudinaryService) Upload(ctx context.Context, fh *multipart.FileHeader) (*UploadResult, error) {
	if fh == nil {
		return nil, ErrImageRequired
	}
	if fh.Size == 0 {
		return nil, ErrEmptyFile
	}
	if fh.Size > MaxImageBytes {
		return nil, ErrImageTooLarge
	}

	src, err := fh.Open()
	if err != nil {
		return nil, ErrInvalidImageFile.WithCause(err)
	}
	defer src.Close()

	data, err := io.ReadAll(io.LimitReader(src, MaxImageBytes+1))
	if err != nil {
		return nil, ErrInvalidImageFile.WithCause(err)
	}
	if len(data) == 0 {
		return nil, ErrEmptyFile
	}
	if len(data) > MaxImageBytes {
		return nil, ErrImageTooLarge
	}

	ext, ok := detectImageFormat(data)
	if !ok {
		return nil, ErrUnsupportedFormat
	}

	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, ErrInvalidImageFile.WithCause(err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > maxImageWidth || cfg.Height > maxImageHeight {
		return nil, ErrImageDimensions
	}

	randomID, err := randomHex(16)
	if err != nil {
		return nil, ErrUploadFailed.WithCause(err)
	}
	publicID := uploadFolder + "/" + randomID
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	signed := map[string]string{
		"allowed_formats": "jpg,png,webp",
		"overwrite":       "false",
		"public_id":       publicID,
		"timestamp":       timestamp,
	}

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	for k, v := range signed {
		_ = w.WriteField(k, v)
	}
	_ = w.WriteField("api_key", s.apiKey)
	_ = w.WriteField("signature", s.sign(signed))

	part, err := w.CreateFormFile("file", randomID+"."+ext)
	if err != nil {
		return nil, ErrUploadFailed.WithCause(err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, ErrUploadFailed.WithCause(err)
	}
	if err := w.Close(); err != nil {
		return nil, ErrUploadFailed.WithCause(err)
	}

	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", s.cloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, &buf)
	if err != nil {
		return nil, ErrUploadFailed.WithCause(err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	body, err := s.call(req, ErrUploadFailed)
	if err != nil {
		return nil, err
	}

	var result cloudinaryUploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, ErrUploadFailed.WithCause(fmt.Errorf("decode response: %w", err))
	}
	if result.Error != nil {
		return nil, ErrStorageReject.WithCause(errors.New(result.Error.Message))
	}
	if result.SecureURL == "" || result.PublicID == "" {
		return nil, ErrUploadFailed.WithCause(errors.New("empty secure_url/public_id in response"))
	}

	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
		Format:   result.Format,
		Width:    result.Width,
		Height:   result.Height,
		Bytes:    result.Bytes,
	}, nil
}

func (s *cloudinaryService) Delete(ctx context.Context, publicID string) error {
	if !publicIDPattern.MatchString(publicID) {
		return ErrInvalidPublicID
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	signed := map[string]string{
		"public_id": publicID,
		"timestamp": timestamp,
	}

	form := url.Values{}
	for k, v := range signed {
		form.Set(k, v)
	}
	form.Set("api_key", s.apiKey)
	form.Set("signature", s.sign(signed))

	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/destroy", s.cloudName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return ErrDeleteFailed.WithCause(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	body, err := s.call(req, ErrDeleteFailed)
	if err != nil {
		return err
	}

	var result cloudinaryDeleteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return ErrDeleteFailed.WithCause(fmt.Errorf("decode response: %w", err))
	}
	if result.Error != nil {
		return ErrDeleteFailed.WithCause(errors.New(result.Error.Message))
	}

	switch result.Result {
	case "ok":
		return nil
	case "not found":
		return ErrImageNotFound
	default:
		return ErrDeleteFailed.WithCause(fmt.Errorf("unexpected result %q", result.Result))
	}
}

// ---------------------------------------------------------------------------
// HTTP call + error classification
// ---------------------------------------------------------------------------

func (s *cloudinaryService) call(req *http.Request, fallback *UploadError) ([]byte, error) {
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, classifyTransportError(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxRespBytes))
	if err != nil {
		return nil, classifyTransportError(err)
	}

	if ue := classifyResponse(resp.StatusCode, resp.Header.Get("Content-Type"), body, fallback); ue != nil {
		return nil, ue
	}
	return body, nil
}

func classifyTransportError(err error) *UploadError {
	var (
		unknownAuth x509.UnknownAuthorityError
		certVerify  *tls.CertificateVerificationError
		dnsErr      *net.DNSError
		netErr      net.Error
	)

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return ErrUploadTimeout.WithCause(err)
	case errors.As(err, &unknownAuth), errors.As(err, &certVerify):
		return ErrUploadBlocked.WithCause(err)
	case errors.As(err, &dnsErr):
		return ErrUploadNetwork.WithCause(err)
	case errors.As(err, &netErr) && netErr.Timeout():
		return ErrUploadTimeout.WithCause(err)
	default:
		return ErrUploadNetwork.WithCause(err)
	}
}

func classifyResponse(status int, contentType string, body []byte, fallback *UploadError) *UploadError {
	trimmed := bytes.TrimSpace(body)
	isJSON := strings.Contains(strings.ToLower(contentType), "json")

	if !isJSON || (len(trimmed) > 0 && trimmed[0] == '<') || status == http.StatusProxyAuthRequired {
		return ErrUploadBlocked.WithCause(fmt.Errorf("non-JSON response: status=%d content-type=%q body=%q",
			status, contentType, snippet(trimmed, 200)))
	}

	cause := fmt.Errorf("status=%d body=%q", status, snippet(trimmed, 300))
	switch {
	case status == http.StatusOK:
		return nil
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return ErrStorageAuth.WithCause(cause)
	case status == 420 || status == http.StatusTooManyRequests:
		return ErrStorageLimited.WithCause(cause)
	case status == http.StatusBadRequest || status == http.StatusUnprocessableEntity || status == http.StatusConflict:
		return ErrStorageReject.WithCause(cause)
	default:
		return fallback.WithCause(cause)
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func (s *cloudinaryService) sign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k, v := range params {
		if v != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+"="+params[k])
	}
	sum := sha1.Sum([]byte(strings.Join(parts, "&") + s.apiSecret))
	return hex.EncodeToString(sum[:])
}

func detectImageFormat(b []byte) (string, bool) {
	switch {
	case len(b) >= 8 && bytes.Equal(b[:8], []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}):
		return "png", true
	case len(b) >= 3 && b[0] == 0xFF && b[1] == 0xD8 && b[2] == 0xFF:
		return "jpg", true
	case len(b) >= 12 && string(b[0:4]) == "RIFF" && string(b[8:12]) == "WEBP":
		return "webp", true
	}
	return "", false
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func snippet(b []byte, max int) string {
	if len(b) > max {
		return string(b[:max]) + "..."
	}
	return string(b)
}
