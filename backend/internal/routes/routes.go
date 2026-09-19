package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
)

func Setup(r *gin.Engine, h *handlers.Handler) {
	api := r.Group("/api")
	api.POST("/login", h.PostLogin)
}