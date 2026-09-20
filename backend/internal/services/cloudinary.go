package services

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/config"
)

var (
	ErrInvalidImageFile  = errors.New("invalid image file")
	ErrImageTooLarge     = errors.New("image file too large, max 1MB")
	ErrUnsupportedFormat = errors.New("unsupported image format, only PNG, JPEG, WEBP allowed")
	ErrImageDimensions   = errors.New("image dimensions too large, max 4000x4000")
	ErrUploadFailed      = errors.New("failed to upload image to cloud storage")
	ErrDeleteFailed      = errors.New("failed to delete image from cloud storage")
)

type UploadResult struct {
	URL      string
	PublicID string
}

type CloudinaryService interface {
	Upload(file *multipart.FileHeader) (*UploadResult, error)
	Delete(publicID string) error
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

func (s *cloudinaryService) Upload(file *multipart.FileHeader) (*UploadResult, error) {
	if file.Size > 1*1024*1024 {
		return nil, ErrImageTooLarge
	}

	src, err := file.Open()
	if err != nil {
		return nil, ErrInvalidImageFile
	}
	defer src.Close()

	header := make([]byte, 26)
	n, err := src.Read(header)
	if err != nil || n < 26 {
		return nil, ErrInvalidImageFile
	}

	if !isValidImage(header) {
		return nil, ErrUnsupportedFormat
	}

	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, ErrInvalidImageFile
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return nil, ErrUploadFailed
	}

	if _, err := io.Copy(part, src); err != nil {
		return nil, ErrUploadFailed
	}

	timestamp := time.Now().Unix()
	signature := s.generateSignature(timestamp)

	_ = writer.WriteField("api_key", s.apiKey)
	_ = writer.WriteField("timestamp", strconv.FormatInt(timestamp, 10))
	_ = writer.WriteField("signature", signature)

	if err := writer.Close(); err != nil {
		return nil, ErrUploadFailed
	}

	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/upload", s.cloudName)

	req, err := http.NewRequest("POST", apiURL, &buf)
	if err != nil {
		return nil, ErrUploadFailed
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, ErrUploadFailed
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrUploadFailed
	}

	var result cloudinaryUploadResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, ErrUploadFailed
	}

	if result.Error != nil {
		return nil, ErrUploadFailed
	}

	if result.SecureURL == "" || result.PublicID == "" {
		return nil, ErrUploadFailed
	}

	return &UploadResult{
		URL:      result.SecureURL,
		PublicID: result.PublicID,
	}, nil
}

func (s *cloudinaryService) Delete(publicID string) error {
	timestamp := time.Now().Unix()
	strToSign := fmt.Sprintf("public_id=%s&timestamp=%d%s", publicID, timestamp, s.apiSecret)
	h := sha1.Sum([]byte(strToSign))
	signature := hex.EncodeToString(h[:])

	form := url.Values{}
	form.Set("public_id", publicID)
	form.Set("api_key", s.apiKey)
	form.Set("timestamp", strconv.FormatInt(timestamp, 10))
	form.Set("signature", signature)

	apiURL := fmt.Sprintf("https://api.cloudinary.com/v1_1/%s/image/destroy", s.cloudName)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return ErrDeleteFailed
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		return ErrDeleteFailed
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ErrDeleteFailed
	}

	var result cloudinaryDeleteResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return ErrDeleteFailed
	}

	if result.Error != nil {
		return ErrDeleteFailed
	}

	if result.Result != "ok" {
		return ErrDeleteFailed
	}

	return nil
}

func (s *cloudinaryService) generateSignature(timestamp int64) string {
	str := fmt.Sprintf("timestamp=%d%s", timestamp, s.apiSecret)
	h := sha1.Sum([]byte(str))
	return hex.EncodeToString(h[:])
}

func isValidImage(header []byte) bool {
	if len(header) < 26 {
		return false
	}

	if header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 &&
		header[4] == 0x0D && header[5] == 0x0A && header[6] == 0x1A && header[7] == 0x0A {
		return true
	}

	if header[0] == 0xFF && header[1] == 0xD8 && header[2] == 0xFF {
		return true
	}

	if len(header) >= 12 &&
		header[0] == 0x52 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x46 &&
		header[8] == 0x57 && header[9] == 0x45 && header[10] == 0x42 && header[11] == 0x50 {
		return true
	}

	return false
}
