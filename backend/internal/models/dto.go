package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type APIResponse struct {
	Success    bool                `json:"success"`
	Message    string              `json:"message"`
	Data       any                 `json:"data,omitempty"`
	Pagination *PaginationResponse `json:"pagination,omitempty"`
}

type PaginationResponse struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
	Token    string `json:"token"`
}

type AuctionRequest struct {
	AuctionName   string          `json:"auctionName" binding:"required"`
	Description   string          `json:"description" binding:"required"`
	ImageLink     string          `json:"imageLink" binding:"required"`
	StartingPrice decimal.Decimal `json:"startingPrice" binding:"required"`
	EndTime       time.Time       `json:"endTime" binding:"required"`
}

type BidWinnerResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type AuctionResponse struct {
	ID            int64              `json:"id"`
	AuctionName   string             `json:"auctionName"`
	Description   string             `json:"description"`
	ImageLink     string             `json:"imageLink"`
	BidWinnerID   *int64             `json:"bidWinnerId"`
	StartingPrice decimal.Decimal    `json:"startingPrice"`
	LastPrice     decimal.Decimal    `json:"lastPrice"`
	CreatedAt     time.Time          `json:"createdAt"`
	EndTime       time.Time          `json:"endTime"`
	IsCompleted   bool               `json:"isCompleted"`
	BidWinner     *BidWinnerResponse `json:"bidWinner"`
}

type GetAuctionsQuery struct {
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	IsCompleted *bool  `form:"isCompleted"`
	Search      string `form:"search"`
}

type BidRequest struct {
	AuctionID int64           `json:"auctionId" binding:"required"`
	BidPrice  decimal.Decimal `json:"bidPrice" binding:"required"`
}

type BidResponse struct {
	ID        int64           `json:"id"`
	AuctionID int64           `json:"auctionId"`
	UserID    int64           `json:"userId"`
	BidPrice  decimal.Decimal `json:"bidPrice"`
	CreatedAt time.Time       `json:"createdAt"`
}

type TopBidResponse struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"userId"`
	UserName  string          `json:"userName"`
	BidPrice  decimal.Decimal `json:"bidPrice"`
	CreatedAt time.Time       `json:"createdAt"`
}

type AuctionDetailResponse struct {
	ID            int64              `json:"id"`
	AuctionName   string             `json:"auctionName"`
	Description   string             `json:"description"`
	ImageLink     string             `json:"imageLink"`
	StartingPrice decimal.Decimal    `json:"startingPrice"`
	LastPrice     decimal.Decimal    `json:"lastPrice"`
	CreatedAt     time.Time          `json:"createdAt"`
	EndTime       time.Time          `json:"endTime"`
	IsCompleted   bool               `json:"isCompleted"`
	BidWinner     *BidWinnerResponse `json:"bidWinner"`
	TotalBids     int64              `json:"totalBids"`
	TotalBidders  int64              `json:"totalBidders"`
	TopBids       []TopBidResponse   `json:"topBids"`
}
