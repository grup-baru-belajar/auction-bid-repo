package handlers

import (
	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
	wsh "github.com/grup-baru-belajar/auction-bid-repo/internal/websocket"
)

type Handler struct {
	authService          services.AuthService
	auctionService       services.AuctionService
	auctionDetailService services.AuctionDetailService
	bidService           services.BidService
	reportingService     services.ReportingService
	wsHandler            *wsh.Handler
	cloudinaryService    services.CloudinaryService
}

func New(
	authService services.AuthService,
	auctionService services.AuctionService,
	auctionDetailService services.AuctionDetailService,
	bidService services.BidService,
	reportingService services.ReportingService,
	wsHandler *wsh.Handler,
	cloudinaryService services.CloudinaryService,
) *Handler {
	return &Handler{
		authService:          authService,
		auctionService:       auctionService,
		auctionDetailService: auctionDetailService,
		bidService:           bidService,
		reportingService:     reportingService,
		wsHandler:            wsHandler,
		cloudinaryService:    cloudinaryService,
	}
}

func respondSuccess(c *gin.Context, status int, message string, data any) {
	c.JSON(status, models.APIResponse{Success: true, Message: message, Data: data})
}

func respondSuccessWithPagination(c *gin.Context, status int, message string, data any, pagination models.PaginationResponse) {
	c.JSON(status, models.APIResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: &pagination,
	})
}

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, models.APIResponse{Success: false, Message: message})
}
