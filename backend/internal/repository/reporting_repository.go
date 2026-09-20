package repository

import (
	"context"
	"database/sql"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

type ReportingRepository interface {
	GetTopAuction(ctx context.Context, limit int) ([]models.TopAuction, error)
	GetAuctionActivity(ctx context.Context) ([]models.AuctionActivity, error)
	GetAuctionStatus(ctx context.Context) ([]models.AuctionStatus, error)
	GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error)
}

type reportingRepository struct {
	db *sql.DB
}

func NewReportingRepository(db *sql.DB) ReportingRepository {
	return &reportingRepository{db: db}
}

func (r *reportingRepository) GetTopAuction(ctx context.Context, limit int) ([]models.TopAuction, error) {
	query := `
	SELECT
		a.id,
		a.auction_name,
		a.starting_price,
		a.last_price AS highest_bid,
		COUNT(b.id) AS total_bids,
		COUNT(DISTINCT b.user_id) AS total_bidders,
		CASE
			WHEN a.is_completed = TRUE
				 OR a.end_time <= CURRENT_TIMESTAMP
			THEN 'ENDED'
			ELSE 'ACTIVE'
		END AS status
	FROM auctions a
	LEFT JOIN bids b
		ON b.auction_id = a.id
	GROUP BY
		a.id,
		a.auction_name,
		a.starting_price,
		a.last_price,
		a.is_completed,
		a.end_time
	ORDER BY total_bids DESC
	LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auctions []models.TopAuction

	for rows.Next() {
		var auction models.TopAuction

		err := rows.Scan(
			&auction.ID,
			&auction.AuctionName,
			&auction.StartingPrice,
			&auction.HighestBid,
			&auction.TotalBids,
			&auction.Bidders,
			&auction.Status,
		)
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, auction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return auctions, nil
}

func (r *reportingRepository) GetAuctionActivity(ctx context.Context) ([]models.AuctionActivity, error) {

	query := `
		SELECT
			DATE(created_at) AS date,
			COUNT(*) AS total_bids
		FROM bids
		WHERE created_at >= CURRENT_DATE - INTERVAL '6 days'
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.AuctionActivity

	for rows.Next() {
		var activity models.AuctionActivity

		err := rows.Scan(
			&activity.Date,
			&activity.TotalBids,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}

func (r *reportingRepository) GetAuctionStatus(ctx context.Context) ([]models.AuctionStatus, error) {
	query := `
		SELECT
			CASE
				WHEN is_completed = TRUE
					OR end_time <= CURRENT_TIMESTAMP
				THEN 'ENDED'
				ELSE 'ACTIVE'
			END AS status,
			COUNT(*) AS total
		FROM auctions
		GROUP BY status
		ORDER BY status
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.AuctionStatus

	for rows.Next() {
		var status models.AuctionStatus

		err := rows.Scan(
			&status.Status,
			&status.Total,
		)
		if err != nil {
			return nil, err
		}

		statuses = append(statuses, status)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return statuses, nil
}



/*
1. get total bidders from every aution and 1 auction in the last x date
*/
func (r *reportingRepository) GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error) {
	query := `
		SELECT COUNT(DISTINCT user_id) AS total_bidders
		FROM bids
		WHERE created_at::DATE >= CURRENT_DATE - $1::INTERVAL
	`
	if auctionId != "" {
		query += ` AND auction_id = $2`
	}
	query += ";"
	var rows *sql.Rows
	var err error
	if auctionId != "" {
		rows, err = r.db.QueryContext(ctx, query, interval, auctionId)
	} else {
		rows, err = r.db.QueryContext(ctx, query, interval)
	}
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var totalBidders int
	if rows.Next() {
		err := rows.Scan(&totalBidders)
		if err != nil {
			return 0, err
		}
	}
	
	if err := rows.Err(); err != nil {
		return 0, err
	}

	return totalBidders, nil
}