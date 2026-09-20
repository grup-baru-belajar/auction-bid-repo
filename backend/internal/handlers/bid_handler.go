package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/middlewares"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/token"
)

func (h *Handler) PostBid(c *gin.Context) {
	var req models.BidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}

	val, exists := c.Get(middlewares.ClaimsKey)
	if !exists {
		respondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	claims, ok := val.(*token.Claims)
	if !ok {
		respondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID, err := claims.UserID()
	if err != nil || userID <= 0 {
		respondError(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.bidService.PlaceBid(c.Request.Context(), userID, req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	go h.wsHandler.BroadcastTopBids(context.Background(), req.AuctionID)

	respondSuccess(c, http.StatusCreated, "Bid placed successfully", result)
}
