package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

func (h *Handler) PostAuction(c *gin.Context) {
	var req models.AuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBindError(c, err)
		return
	}

	result, err := h.auctionService.CreateAuction(c.Request.Context(), req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	respondSuccess(c, http.StatusCreated, "Auction created successfully", result)
}

func (h *Handler) GetAuctions(c *gin.Context) {
	var query models.GetAuctionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		respondBindError(c, err)
		return
	}

	auctions, pagination, err := h.auctionService.GetAuctions(c.Request.Context(), query)
	if err != nil {
		respondWithError(c, err)
		return
	}

	respondSuccessWithPagination(c, http.StatusOK, "Auctions retrieved successfully", auctions, pagination)
}
