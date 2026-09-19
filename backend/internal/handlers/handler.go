package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
)

type Handler struct {
	authService services.AuthService
}

func New(authService services.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}

func respondSuccess(c *gin.Context, status int, message string, data any) {
	c.JSON(status, models.APIResponse{Success: true, Message: message, Data: data})
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, models.APIResponse{Success: false, Message: message})
}