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
	GetTotalTransaction(ctx context.Context, interval string) (models.TotalTransaction, error)
	GetAuctionSummary(ctx context.Context) (models.AuctionSummary, error)
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
		WHERE created_at >= CURRENT_DATE - CAST($1 AS INTERVAL)
	`
	if auctionId != "" {
		query += ` AND auction_id = $2`
	}
	query += ";"
	var row *sql.Row
	var err error
	if auctionId != "" {
		row = r.db.QueryRowContext(ctx, query, interval, auctionId)
	} else {
		row = r.db.QueryRowContext(ctx, query, interval)
	}
	if err != nil {
		return 0, err
	}

	var totalBidders int
	if err := row.Scan(&totalBidders); err != nil {
		return 0, err
	}

	return totalBidders, nil
}

/*
- count total transaction happen in an interval of time
- sum total gross sales in an interval of time
*/
func (r *reportingRepository) GetTotalTransaction(ctx context.Context, interval string) (models.TotalTransaction, error) {
	query := `
		SELECT COUNT(*) AS transaction_count, COALESCE(SUM(last_price), 0) AS total_gross_sales
		FROM auctions
		WHERE created_at >= CURRENT_DATE - CAST($1 AS INTERVAL) AND is_completed = TRUE AND end_time <= CURRENT_TIMESTAMP;
	`
	row := r.db.QueryRowContext(ctx, query, interval)

	var totalTransactions models.TotalTransaction
	if err := row.Scan(&totalTransactions.TransactionCount, &totalTransactions.TotalGrossSales); err != nil {
		return models.TotalTransaction{}, err
	}

	return totalTransactions, nil
}

func (r *reportingRepository) GetAuctionSummary(ctx context.Context) (models.AuctionSummary, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM auctions) AS total_auctions,

			(SELECT COUNT(*) FROM auctions
			 WHERE is_completed IS DISTINCT FROM TRUE
			   AND end_time > NOW()) AS ongoing_auctions,

			(SELECT COUNT(*) FROM auctions
			 WHERE is_completed IS TRUE
				OR end_time <= NOW()) AS completed_auctions,

			(SELECT COUNT(b.id) FROM auctions a
			 JOIN bids b ON b.auction_id = a.id
			 WHERE a.is_completed IS DISTINCT FROM TRUE
			   AND a.end_time > NOW()) AS total_bids_ongoing,

			(SELECT COUNT(b.id) FROM auctions a
			 JOIN bids b ON b.auction_id = a.id
			 WHERE a.is_completed IS TRUE
				OR a.end_time <= NOW()) AS total_bids_completed,

			(SELECT COUNT(*) FROM bids) AS total_bids_all
	`

	var s models.AuctionSummary
	err := r.db.QueryRowContext(ctx, query).Scan(
		&s.TotalAuctions,
		&s.OngoingAuctions,
		&s.CompletedAuctions,
		&s.TotalBidsOngoing,
		&s.TotalBidsCompleted,
		&s.TotalBidsAll,
	)
	if err != nil {
		return s, err
	}

	return s, nil
}
