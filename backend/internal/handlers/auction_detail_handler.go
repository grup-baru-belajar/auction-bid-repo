package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAuctionByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "Invalid auction id")
		return
	}

	result, err := h.auctionDetailService.GetAuctionByID(c.Request.Context(), id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	respondSuccess(c, http.StatusOK, "Auction retrieved successfully", result)
}
