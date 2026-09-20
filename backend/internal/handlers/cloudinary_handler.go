package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostUploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		log.Printf("[Upload] No file received: %v", err)
		respondError(c, http.StatusBadRequest, "Image file is required")
		return
	}

	log.Printf("[Upload] File received: name=%s, size=%d, header=%v", file.Filename, file.Size, file.Header)

	result, err := h.cloudinaryService.Upload(file)
	if err != nil {
		log.Printf("[Upload] Upload failed: %v", err)
		respondError(c, http.StatusBadRequest, "Invalid image file")
		return
	}

	log.Printf("[Upload] Upload success: url=%s", result.URL)
	respondSuccess(c, http.StatusOK, "Image uploaded successfully", gin.H{
		"url":       result.URL,
		"public_id": result.PublicID,
	})
}

func (h *Handler) DeleteUploadImage(c *gin.Context) {
	publicID := c.Query("public_id")
	if publicID == "" {
		respondError(c, http.StatusBadRequest, "public_id is required")
		return
	}

	log.Printf("[Upload] Deleting image: public_id=%s", publicID)

	if err := h.cloudinaryService.Delete(publicID); err != nil {
		log.Printf("[Upload] Delete failed: %v", err)
		respondError(c, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	log.Printf("[Upload] Delete success: public_id=%s", publicID)
	respondSuccess(c, http.StatusOK, "Image deleted successfully", nil)
}
