package models

import "github.com/shopspring/decimal"

type TopAuction struct {
	ID            int64   `json:"id"`
	AuctionName   string  `json:"auctionName"`
	StartingPrice float64 `json:"startingPrice"`
	HighestBid    float64 `json:"highestBid"`
	TotalBids     int64   `json:"totalBid"`
	Bidders       int64   `json:"bidders"`
	Status        string  `json:"status"`
}

// AuctionActivity is auctions created and bids placed, per day.
type AuctionActivity struct {
	Date          string `json:"date"`
	TotalAuctions int64  `json:"totalAuctions"`
	TotalBids     int64  `json:"totalBids"`
}

// TransactionWeek is the sum of completed auctions' lastPrice for one week.
type TransactionWeek struct {
	WeekStart string  `json:"weekStart"`
	Total     float64 `json:"total"`
}

type AuctionStatus struct {
	Status string `json:"status"`
	Total  int64  `json:"total"`
}

type TotalTransaction struct {
	TransactionCount string `json:"transactionCount"`
	TotalGrossSales decimal.Decimal  `json:"totalGrossSales"`
}

// AuctionSummary is the aggregate snapshot returned by GET /reporting/auction-summary
// and broadcast to reporting WebSocket clients after each bid.
type AuctionSummary struct {
	TotalAuctions      int64 `json:"totalAuctions"`
	OngoingAuctions    int64 `json:"ongoingAuctions"`
	CompletedAuctions  int64 `json:"completedAuctions"`
	TotalBidsOngoing   int64 `json:"totalBidsOngoing"`
	TotalBidsCompleted int64 `json:"totalBidsCompleted"`
	TotalBidsAll       int64 `json:"totalBidsAll"`
}
type TopSpenderBidder struct {
	UserID          int64   `json:"user_id"`
	Username        string  `json:"username"`
	Name            string  `json:"name"`
	TotalMoneySpent decimal.Decimal `json:"total_money_spent"`
	LastBidAgo      string  `json:"last_bid_ago"`
}