package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) PostUploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		respondError(c, http.StatusBadRequest, "Image file is required")
		return
	}

	result, err := h.cloudinaryService.Upload(file)
	if err != nil {
		respondError(c, http.StatusBadRequest, "Invalid image file")
		return
	}

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

	if err := h.cloudinaryService.Delete(publicID); err != nil {
		log.Printf("[Upload] Delete failed: %v", err)
		respondError(c, http.StatusInternalServerError, "Failed to delete image")
		return
	}

	respondSuccess(c, http.StatusOK, "Image deleted successfully", nil)
}
