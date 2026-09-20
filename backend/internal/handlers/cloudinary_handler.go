package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
)

const maxUploadRequestBytes = services.MaxImageBytes + 64*1024

var uploadErrorStatus = map[string]int{
	"INVALID_REQUEST":      http.StatusBadRequest,
	"IMAGE_REQUIRED":       http.StatusBadRequest,
	"EMPTY_FILE":           http.StatusBadRequest,
	"INVALID_PUBLIC_ID":    http.StatusBadRequest,
	"FILE_TOO_LARGE":       http.StatusRequestEntityTooLarge,
	"UNSUPPORTED_FORMAT":   http.StatusUnsupportedMediaType,
	"INVALID_IMAGE":        http.StatusUnprocessableEntity,
	"IMAGE_DIMENSIONS":     http.StatusUnprocessableEntity,
	"IMAGE_NOT_FOUND":      http.StatusNotFound,
	"STORAGE_REJECTED":     http.StatusUnprocessableEntity,
	"STORAGE_RATE_LIMITED": http.StatusTooManyRequests,
	"NETWORK_BLOCKED":      http.StatusBadGateway,
	"NETWORK_ERROR":        http.StatusBadGateway,
	"UPLOAD_TIMEOUT":       http.StatusGatewayTimeout,
	"UPLOAD_FAILED":        http.StatusBadGateway,
	"DELETE_FAILED":        http.StatusBadGateway,
	"STORAGE_CONFIG_ERROR": http.StatusInternalServerError,
}

func (h *Handler) PostUploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadRequestBytes)

	file, err := c.FormFile("image")
	if err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxErr):
			respondUploadError(c, services.ErrImageTooLarge)
		case errors.Is(err, http.ErrMissingFile):
			respondUploadError(c, services.ErrImageRequired)
		default:
			respondUploadError(c, services.ErrInvalidRequest.WithCause(err))
		}
		return
	}

	result, err := h.cloudinaryService.Upload(c.Request.Context(), file)
	if err != nil {
		respondUploadError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, "Image uploaded successfully", gin.H{
		"url":       result.URL,
		"public_id": result.PublicID,
		"format":    result.Format,
		"width":     result.Width,
		"height":    result.Height,
		"bytes":     result.Bytes,
	})
}

func (h *Handler) DeleteUploadImage(c *gin.Context) {
	publicID := c.Query("public_id")
	if publicID == "" {
		respondUploadError(c, services.ErrInvalidPublicID)
		return
	}

	if err := h.cloudinaryService.Delete(c.Request.Context(), publicID); err != nil {
		respondUploadError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, "Image deleted successfully", nil)
}

func respondUploadError(c *gin.Context, err error) {
	var ue *services.UploadError
	if !errors.As(err, &ue) {
		log.Printf("[Upload] unexpected error: %v", err)
		ue = services.ErrUploadFailed
	} else if ue.Err != nil {
		log.Printf("[Upload] %s (%s %s): %v", ue.Code, c.Request.Method, c.FullPath(), ue.Err)
	}

	status, ok := uploadErrorStatus[ue.Code]
	if !ok {
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{
		"success": false,
		"message": ue.Message,
		"error": gin.H{
			"code": ue.Code,
		},
	})
}
