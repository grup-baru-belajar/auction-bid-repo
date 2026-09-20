package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/shopspring/decimal"
)

var (
	ErrAuctionCompleted = errors.New("auction already completed")
	ErrBidTooLow        = errors.New("bid price must be greater than current price")
)

type BidRepository struct {
	db *sql.DB
}

func NewBidRepository(db *sql.DB) *BidRepository {
	return &BidRepository{db: db}
}

func (r *BidRepository) Place(ctx context.Context, auctionID, userID int64, bidPrice decimal.Decimal) (*models.Bid, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin bid tx: %w", err)
	}
	defer tx.Rollback()

	var lastPrice decimal.Decimal
	var endTime time.Time
	var isCompleted bool

	err = tx.QueryRowContext(ctx, `
		SELECT last_price, end_time, is_completed
		FROM auctions
		WHERE id = $1
		FOR UPDATE
	`, auctionID).Scan(&lastPrice, &endTime, &isCompleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuctionNotFound
		}
		return nil, fmt.Errorf("lock auction for bid: %w", err)
	}

	if isCompleted || !endTime.After(time.Now()) {
		return nil, ErrAuctionCompleted
	}

	if !bidPrice.GreaterThan(lastPrice) {
		return nil, ErrBidTooLow
	}

	bid := &models.Bid{
		AuctionID: auctionID,
		UserID:    userID,
		BidPrice:  bidPrice,
	}

	err = tx.QueryRowContext(ctx, `
		INSERT INTO bids (auction_id, user_id, bid_price)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, auctionID, userID, bidPrice).Scan(&bid.ID, &bid.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert bid: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE auctions
		SET last_price = $1, bid_winner_id = $2
		WHERE id = $3
	`, bidPrice, userID, auctionID)
	if err != nil {
		return nil, fmt.Errorf("update auction after bid: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit bid tx: %w", err)
	}

	return bid, nil
}
