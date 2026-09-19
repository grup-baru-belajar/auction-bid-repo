package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTopAuction(c *gin.Context) {
	auctions, err := h.reportingService.GetTopAuction(c, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get top auctions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Top auctions retrieved successfully",
		"data":    auctions,
	})
}

func (h *Handler) GetAuctionActivity(c *gin.Context) {
	activities, err := h.reportingService.GetAuctionActivity(c.Request.Context())

	if err != nil {
		log.Printf("GetAuctionActivity error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get auction activity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Auction activity retrieved successfully",
		"data":    activities,
	})
}

func (h *Handler) GetAuctionStatus(c *gin.Context) {
	statuses, err := h.reportingService.GetAuctionStatus(c.Request.Context())

	if err != nil {
		log.Printf("GetAuctionStatus error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get auction status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Auction status retrieved successfully",
		"data":    statuses,
	})
}
