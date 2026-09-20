package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

type AuctionWithWinner struct {
	Auction       models.Auction
	BidWinnerName *string
}

type AuctionRepository struct {
	db *sql.DB
}

func NewAuctionRepository(db *sql.DB) *AuctionRepository {
	return &AuctionRepository{db: db}
}

func (r *AuctionRepository) Create(ctx context.Context, auction *models.Auction) (*models.Auction, error) {
	query := `
		INSERT INTO auctions (
			auction_name, description, image_link, starting_price, last_price, end_time, is_completed
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		auction.AuctionName,
		auction.Description,
		auction.ImageLink,
		auction.StartingPrice,
		auction.LastPrice,
		auction.EndTime,
		auction.IsCompleted,
	).Scan(&auction.ID, &auction.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create auction: %w", err)
	}

	return auction, nil
}

func (r *AuctionRepository) FindAll(ctx context.Context, limit, offset int, isCompleted *bool, search string) ([]AuctionWithWinner, int64, error) {
	countQuery := `SELECT COUNT(*) FROM auctions a`
	selectQuery := `
		SELECT 
			a.id, a.auction_name, a.description, a.image_link, a.bid_winner_id,
			a.starting_price, a.last_price, a.created_at, a.end_time,
			(a.is_completed OR a.end_time <= NOW()) AS is_completed,
			u.name AS bid_winner_name
		FROM auctions a
		LEFT JOIN users u ON a.bid_winner_id = u.id
	`

	var conditions []string
	var args []any
	argIndex := 1

	if isCompleted != nil {
		conditions = append(conditions, fmt.Sprintf("(a.is_completed OR a.end_time <= NOW()) = $%d", argIndex))
		args = append(args, *isCompleted)
		argIndex++
	}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("a.auction_name ILIKE $%d", argIndex))
		args = append(args, "%"+search+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		countQuery += whereClause
		selectQuery += whereClause
	}

	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count auctions: %w", err)
	}

	selectQuery += fmt.Sprintf(" ORDER BY a.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	selectArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("find all auctions: %w", err)
	}
	defer rows.Close()

	var result []AuctionWithWinner
	for rows.Next() {
		var item AuctionWithWinner
		var bidWinnerName sql.NullString

		err := rows.Scan(
			&item.Auction.ID,
			&item.Auction.AuctionName,
			&item.Auction.Description,
			&item.Auction.ImageLink,
			&item.Auction.BidWinnerID,
			&item.Auction.StartingPrice,
			&item.Auction.LastPrice,
			&item.Auction.CreatedAt,
			&item.Auction.EndTime,
			&item.Auction.IsCompleted,
			&bidWinnerName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan auction: %w", err)
		}

		if bidWinnerName.Valid {
			item.BidWinnerName = &bidWinnerName.String
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows err: %w", err)
	}

	return result, total, nil
}
