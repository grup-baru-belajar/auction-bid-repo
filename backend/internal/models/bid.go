package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Bid struct {
	ID        int64
	AuctionID int64
	UserID    int64
	BidPrice  decimal.Decimal
	CreatedAt time.Time
}
