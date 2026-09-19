package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

var ErrAuctionNotFound = errors.New("auction not found")

type AuctionDetail struct {
	Auction models.Auction
	BidWinnerName *string
	TotalBids int64
	TotalBidders int64
}

type TopBid struct {
	Bid models.Bid
	UserName string
}

type AuctionDetailRepository interface {
	FindByID(ctx context.Context, id int64) (*AuctionDetail, error)
	FindTopBids(ctx context.Context, auctionID int64, limit int) ([]TopBid, error)
}

type auctionDetailRepository struct {
	db *sql.DB
}

func NewAuctionDetailRepository(db *sql.DB) AuctionDetailRepository {
	return &auctionDetailRepository{db: db}
}

func (r *auctionDetailRepository) FindByID(ctx context.Context, id int64) (*AuctionDetail, error) {
	query := `
		SELECT
			a.id, a.auction_name, a.description, a.image_link, a.bid_winner_id,
			a.starting_price, a.last_price, a.created_at, a.end_time, a.is_completed,
			u.name AS bid_winner_name,
			(SELECT COUNT(*) FROM bids b WHERE b.auction_id = a.id) AS total_bids,
			(SELECT COUNT(DISTINCT b.user_id) FROM bids b WHERE b.auction_id = a.id) AS total_bidders
		FROM auctions a
		LEFT JOIN users u ON u.id = a.bid_winner_id
		WHERE a.id = $1
	`

	var detail AuctionDetail
	var bidWinnerName sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.Auction.ID,
		&detail.Auction.AuctionName,
		&detail.Auction.Description,
		&detail.Auction.ImageLink,
		&detail.Auction.BidWinnerID,
		&detail.Auction.StartingPrice,
		&detail.Auction.LastPrice,
		&detail.Auction.CreatedAt,
		&detail.Auction.EndTime,
		&detail.Auction.IsCompleted,
		&bidWinnerName,
		&detail.TotalBids,
		&detail.TotalBidders,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuctionNotFound
		}
		return nil, fmt.Errorf("find auction by id: %w", err)
	}

	if bidWinnerName.Valid {
		detail.BidWinnerName = &bidWinnerName.String
	}

	return &detail, nil
}

func (r *auctionDetailRepository) FindTopBids(ctx context.Context, auctionID int64, limit int) ([]TopBid, error) {
	query := `
		SELECT b.id, b.auction_id, b.user_id, u.name, b.bid_price, b.created_at
		FROM bids b
		JOIN users u ON u.id = b.user_id
		WHERE b.auction_id = $1
		ORDER BY b.bid_price DESC, b.created_at ASC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, auctionID, limit)
	if err != nil {
		return nil, fmt.Errorf("find top bids: %w", err)
	}
	defer rows.Close()

	result := make([]TopBid, 0, limit)
	for rows.Next() {
		var item TopBid
		err := rows.Scan(
			&item.Bid.ID,
			&item.Bid.AuctionID,
			&item.Bid.UserID,
			&item.UserName,
			&item.Bid.BidPrice,
			&item.Bid.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan top bid: %w", err)
		}
		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, nil
}

