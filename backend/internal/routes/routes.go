package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/middlewares"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/token"
	wsh "github.com/grup-baru-belajar/auction-bid-repo/internal/websocket"
)

func Setup(r *gin.Engine, h *handlers.Handler, tokenManager *token.TokenManager, wsHandler *wsh.Handler) {
	api := r.Group("/api/v1")
	api.POST("/login", h.PostLogin)

	api.GET("/auctions", h.GetAuctions)
	api.GET("/auctions/:id", h.GetAuctionByID)
	api.POST("/auctions", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.PostAuction)

	api.POST("/upload/image", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.PostUploadImage)
	api.DELETE("/upload/image", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.DeleteUploadImage)

	api.POST("/bid", middlewares.Auth(tokenManager), h.PostBid)

	api.GET("/reporting/top-auction", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetTopAuction)
	api.GET("/reporting/auction-activity", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetAuctionActivity)
	api.GET("/reporting/auction-status", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetAuctionStatus)
	api.GET("/reporting/total-bidders", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetTotalBidders)
	api.GET("/reporting/total-transaction", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetTotalTransaction)

	api.GET("/reporting/auction-summary", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetAuctionSummary)
	api.GET("/reporting/transaction-overview", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetTransactionOverview)
	
	// can have sorting based on amount or transaction_count, default is amount
	api.GET("/reporting/top-bidders", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetTopBidders)

	r.GET("/ws/auctions/:id/top-bids", wsHandler.ServeHTTP)
}
