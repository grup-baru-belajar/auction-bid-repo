package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func queryIntDefault(c *gin.Context, key string, def int) int {
	raw := c.Query(key)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v <= 0 {
		return def
	}
	return v
}

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
	interval := queryIntDefault(c, "interval", 7)
	activities, err := h.reportingService.GetAuctionActivity(c.Request.Context(), interval)

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

func (h *Handler) GetTransactionOverview(c *gin.Context) {
	weeks := queryIntDefault(c, "weeks", 4)
	overview, err := h.reportingService.GetTransactionOverview(c.Request.Context(), weeks)

	if err != nil {
		log.Printf("GetTransactionOverview error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get transaction overview",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaction overview retrieved successfully",
		"data":    overview,
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

func (h *Handler) GetTotalBidders(c *gin.Context) {
	interval := c.Query("interval")
	if interval == "" {
		interval = "7"
	}
	if interval != "all" {
		interval = interval + " days"
	}
	auctionId := c.Query("auctionId")

	totalBidders, err := h.reportingService.GetTotalBidders(c.Request.Context(), interval, auctionId)

	if err != nil {
		log.Printf("GetTotalBidders error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get total bidders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Total bidders retrieved successfully",
		"data":    totalBidders,
	})
}

func (h *Handler) GetTotalTransaction(c *gin.Context) {
	interval := c.Query("interval")
	if interval == "" {
		interval = "7"
	}
	interval = interval + " days"
	totalTransactions, err := h.reportingService.GetTotalTransaction(c.Request.Context(), interval)

	if err != nil {
		log.Printf("GetTotalTransaction error: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get total transactions",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Total transactions retrieved successfully",
		"data":    totalTransactions,
	})
}

func (h *Handler) GetAuctionSummary(c *gin.Context) {
	summary, err := h.reportingService.GetAuctionSummary(c.Request.Context())
	if err != nil {
		log.Printf("GetAuctionSummary error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get auction summary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Auction summary retrieved successfully",
		"data":    summary,
	})
}

func (h *Handler) GetTopBidders(c *gin.Context) {
	limit := c.Query("limit")
	if limit == "" {
		limit = "5"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid limit parameter",
		})
		return
	}

	sortBy := c.Query("sort_by")
	if sortBy == "" {
		sortBy = "amount"
	}

	topBidders, err := h.reportingService.GetTopBidders(c.Request.Context(), limitInt, sortBy)
	if err != nil {
		log.Printf("GetTopBidders error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get top bidders data",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Top bidders by money spent retrieved successfully",
		"data":    topBidders,
	})
}
