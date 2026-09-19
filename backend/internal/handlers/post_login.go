package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
)

func (h *Handler) PostLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Username and password are required")
		return
	}

	res, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			respondError(c, http.StatusUnauthorized, "Invalid username or password")
			return
		}
		log.Printf("login failes : %v", err)
		respondError(c, http.StatusInternalServerError, "Internal server error")
		return
	}

	respondSuccess(c, http.StatusOK, "Login successful", res)
}