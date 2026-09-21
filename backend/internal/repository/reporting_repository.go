package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

type ReportingRepository interface {
	GetTopAuction(ctx context.Context, limit int) ([]models.TopAuction, error)
	GetAuctionActivity(ctx context.Context, intervalDays int) ([]models.AuctionActivity, error)
	GetAuctionStatus(ctx context.Context) ([]models.AuctionStatus, error)
	GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error)
	GetTotalTransaction(ctx context.Context, interval string) (models.TotalTransaction, error)
	GetAuctionSummary(ctx context.Context) (models.AuctionSummary, error)
	GetTransactionOverview(ctx context.Context, weeks int) ([]models.TransactionWeek, error)
	GetTopBidders(ctx context.Context, limit int, sortBy string) ([]models.TopBidder, error)
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

// GetAuctionActivity returns, per day, how many auctions were created and
// how many bids were placed — one row per date that has either.
func (r *reportingRepository) GetAuctionActivity(ctx context.Context, intervalDays int) ([]models.AuctionActivity, error) {

	query := `
		SELECT
			COALESCE(a.date, b.date) AS date,
			COALESCE(a.total_auctions, 0) AS total_auctions,
			COALESCE(b.total_bids, 0) AS total_bids
		FROM (
			SELECT DATE(created_at) AS date, COUNT(*) AS total_auctions
			FROM auctions
			WHERE created_at >= CURRENT_DATE - ($1 || ' days')::INTERVAL
			GROUP BY DATE(created_at)
		) a
		FULL OUTER JOIN (
			SELECT DATE(created_at) AS date, COUNT(*) AS total_bids
			FROM bids
			WHERE created_at >= CURRENT_DATE - ($1 || ' days')::INTERVAL
			GROUP BY DATE(created_at)
		) b ON a.date = b.date
		ORDER BY date ASC
	`

	rows, err := r.db.QueryContext(ctx, query, intervalDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.AuctionActivity

	for rows.Next() {
		var activity models.AuctionActivity

		err := rows.Scan(
			&activity.Date,
			&activity.TotalAuctions,
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
1. get total bidders from every auction and 1 auction in the last x date
*/
func (r *reportingRepository) GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error) {
	var query string
	args := []any{}
	if interval == "all" {
		query = `SELECT COUNT(DISTINCT user_id) AS total_bidders FROM bids`
		if auctionId != "" {
			query += ` WHERE auction_id = $1`
			args = append(args, auctionId)
		}
	} else {
		query = `
			SELECT COUNT(DISTINCT user_id) AS total_bidders
			FROM bids
			WHERE created_at::DATE >= CURRENT_DATE - $1::INTERVAL
		`
		args = append(args, interval)
		if auctionId != "" {
			query += ` AND auction_id = $2`
			args = append(args, auctionId)
		}
	}
	query += ";"

	rows, err := r.db.QueryContext(ctx, query, args...)
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

// GetTransactionOverview sums the lastPrice of completed auctions per week,
// for the last `weeks` calendar weeks (weeks start on Monday).
func (r *reportingRepository) GetTransactionOverview(ctx context.Context, weeks int) ([]models.TransactionWeek, error) {
	query := `
		SELECT
			DATE_TRUNC('week', end_time)::DATE AS week_start,
			COALESCE(SUM(last_price), 0) AS total
		FROM auctions
		WHERE (is_completed = TRUE OR end_time <= CURRENT_TIMESTAMP)
			AND end_time >= CURRENT_DATE - ($1 || ' weeks')::INTERVAL
		GROUP BY week_start
		ORDER BY week_start ASC
	`

	rows, err := r.db.QueryContext(ctx, query, weeks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.TransactionWeek

	for rows.Next() {
		var week models.TransactionWeek

		if err := rows.Scan(&week.WeekStart, &week.Total); err != nil {
			return nil, err
		}

		result = append(result, week)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *reportingRepository) GetTopBidders(ctx context.Context, limit int, sortBy string) ([]models.TopBidder, error) {
	query := `
		SELECT 
		u.id, 
		u.username, 
		u.name, 
		COALESCE(SUM(a.last_price), 0) AS total_money_spent, 
		COUNT(a.bid_winner_id) AS total_wins,
		(CURRENT_DATE - 
			(
				SELECT MAX(b.created_at)::DATE 
				FROM bids b 
				WHERE b.user_id = u.id
			)
		) || ' days ago' AS last_bid_ago
		FROM users u 
		INNER JOIN auctions a ON u.id = a.bid_winner_id
		WHERE a.is_completed = TRUE AND a.end_time <= CURRENT_TIMESTAMP
		GROUP BY u.id, u.username, u.name
	`
	if sortBy == "amount" {
		query += ` ORDER BY total_money_spent DESC`
	} else if sortBy == "transaction_count" {
		query += ` ORDER BY total_wins DESC`
	} else {
		return nil, fmt.Errorf("invalid sortBy value: %s", sortBy)
	}
	query += ` LIMIT $1`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topBidders []models.TopBidder

	for rows.Next() {
		var bidder models.TopBidder
		if err := rows.Scan(&bidder.UserID, &bidder.Username, &bidder.Name, &bidder.TotalMoneySpent, &bidder.TotalWins, &bidder.LastBidAgo); err != nil {
			return nil, err
		}
		topBidders = append(topBidders, bidder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topBidders, nil
}
