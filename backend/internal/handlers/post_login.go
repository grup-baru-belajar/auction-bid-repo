package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

func (h *Handler) PostLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, "Login successful", res)
}
