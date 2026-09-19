package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Auction struct {
	ID            int64
	AuctionName   string
	Description   string
	ImageLink     string
	BidWinnerID   *int64
	StartingPrice decimal.Decimal
	LastPrice     decimal.Decimal
	CreatedAt     time.Time
	EndTime       time.Time
	IsCompleted   bool
}