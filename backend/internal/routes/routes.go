package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/middlewares"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
)

func Setup(r *gin.Engine, h *handlers.Handler, tokenManager *services.TokenManager) {
	api := r.Group("/api/v1")
	api.POST("/login", h.PostLogin)

	api.GET("/auctions", h.GetAuctions)
	api.GET("/auctions/:id", h.GetAuctionByID)
	api.POST("/auctions", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.PostAuction)

	api.POST("/bid", middlewares.Auth(tokenManager), h.PostBid)

	api.GET("/reporting/top-auction", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetTopAuction)
	api.GET("/reporting/auction-activity", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetAuctionActivity)
	api.GET("/reporting/auction-status", middlewares.Auth(tokenManager), middlewares.RequireAdmin(), h.GetAuctionStatus)

}
