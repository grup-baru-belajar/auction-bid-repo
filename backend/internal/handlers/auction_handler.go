package handlers

import (
	"log"
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

	log.Printf("[Auction] Received create request: name=%q, price=%s, endTime=%s, imageLink=%s",
		req.AuctionName, req.StartingPrice.String(), req.EndTime.Format("2006-01-02T15:04:05Z"), req.ImageLink)

	result, err := h.auctionService.CreateAuction(c.Request.Context(), req)
	if err != nil {
		log.Printf("[Auction] Create failed: %v", err)
		respondWithError(c, err)
		return
	}

	log.Printf("[Auction] Auction created successfully: id=%d", result.ID)
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
