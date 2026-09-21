//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func TestAuctionDetailRepository_EndToEnd(t *testing.T) {
	dsn := os.Getenv("PG_DSN")
	if dsn == "" {
		t.Skip("PG_DSN not set; skipping integration test")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// Use single connection so SET search_path applies for all statements
	db.SetMaxOpenConns(1)

	schema := fmt.Sprintf("test_schema_%d", time.Now().UnixNano())
	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE SCHEMA %s", schema)); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	// ensure cleanup
	defer func() {
		_, _ = db.ExecContext(context.Background(), fmt.Sprintf("DROP SCHEMA %s CASCADE", schema))
	}()

	if _, err := db.ExecContext(context.Background(), fmt.Sprintf("SET search_path TO %s", schema)); err != nil {
		t.Fatalf("set search_path: %v", err)
	}

	// create minimal tables required by auction detail repository
	// timestamptz mirrors the real migration, which matters for the end_time comparison
	createSQL := `
    CREATE TABLE users (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL
    );

    CREATE TABLE auctions (
        id BIGSERIAL PRIMARY KEY,
        auction_name VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        image_link TEXT NOT NULL,
        bid_winner_id BIGINT NULL,
        starting_price NUMERIC(18,2) NOT NULL,
        last_price NUMERIC(18,2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
        end_time TIMESTAMPTZ NOT NULL,
        is_completed BOOLEAN NOT NULL DEFAULT FALSE
    );

    CREATE TABLE bids (
        id BIGSERIAL PRIMARY KEY,
        auction_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        bid_price NUMERIC(18,2) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
    );
    `

	if _, err := db.ExecContext(context.Background(), createSQL); err != nil {
		t.Fatalf("create tables: %v", err)
	}

	ctx := context.Background()

	insertUser := func(name string) int64 {
		var id int64
		if err := db.QueryRowContext(ctx, `INSERT INTO users (name) VALUES ($1) RETURNING id`, name).Scan(&id); err != nil {
			t.Fatalf("insert user %q: %v", name, err)
		}
		return id
	}

	insertAuction := func(name string, winnerID *int64, endTime time.Time, completed bool) int64 {
		var id int64
		err := db.QueryRowContext(ctx,
			`INSERT INTO auctions (auction_name, description, image_link, bid_winner_id, starting_price, last_price, end_time, is_completed)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`,
			name, "desc", "img", winnerID, decimal.NewFromInt(1000000), decimal.NewFromInt(1000000), endTime, completed).Scan(&id)
		if err != nil {
			t.Fatalf("insert auction %q: %v", name, err)
		}
		return id
	}

	insertBid := func(auctionID, userID int64, price int64, createdAt time.Time) {
		_, err := db.ExecContext(ctx,
			`INSERT INTO bids (auction_id, user_id, bid_price, created_at) VALUES ($1,$2,$3,$4)`,
			auctionID, userID, decimal.NewFromInt(price), createdAt)
		if err != nil {
			t.Fatalf("insert bid: %v", err)
		}
	}

	siti := insertUser("Siti Rahma")
	john := insertUser("John Doe")

	now := time.Now()
	running := insertAuction("running auction", &siti, now.Add(1*time.Hour), false)
	expired := insertAuction("expired auction", nil, now.Add(-1*time.Hour), false)
	flagged := insertAuction("flagged auction", nil, now.Add(1*time.Hour), true)

	// four bids from two bidders, two of them share the same price
	insertBid(running, john, 1000000, now.Add(-3*time.Hour))
	insertBid(running, siti, 2000000, now.Add(-2*time.Hour))
	insertBid(running, john, 2000000, now.Add(-1*time.Hour))
	insertBid(running, siti, 1500000, now.Add(-30*time.Minute))

	repo := NewAuctionDetailRepository(db)

	t.Run("find by id returns joined winner and counts", func(t *testing.T) {
		detail, err := repo.FindByID(ctx, running)
		if err != nil {
			t.Fatalf("find by id: %v", err)
		}
		if detail.Auction.ID != running || detail.Auction.AuctionName != "running auction" {
			t.Fatalf("unexpected auction returned: %+v", detail.Auction)
		}
		if detail.BidWinnerName == nil || *detail.BidWinnerName != "Siti Rahma" {
			t.Fatalf("expected winner name from the join, got %v", detail.BidWinnerName)
		}
		if detail.TotalBids != 4 {
			t.Fatalf("expected 4 bids, got %d", detail.TotalBids)
		}
		if detail.TotalBidders != 2 {
			t.Fatalf("expected 2 distinct bidders, got %d", detail.TotalBidders)
		}
	})

	t.Run("is_completed is derived from end_time", func(t *testing.T) {
		tests := []struct {
			name      string
			auctionID int64
			want      bool
		}{
			{"still running", running, false},
			{"end time already passed", expired, true},
			{"flagged completed while still running", flagged, true},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				detail, err := repo.FindByID(ctx, tc.auctionID)
				if err != nil {
					t.Fatalf("find by id: %v", err)
				}
				if detail.Auction.IsCompleted != tc.want {
					t.Fatalf("expected is_completed %v, got %v", tc.want, detail.Auction.IsCompleted)
				}
			})
		}
	})

	t.Run("no winner leaves the join empty", func(t *testing.T) {
		detail, err := repo.FindByID(ctx, expired)
		if err != nil {
			t.Fatalf("find by id: %v", err)
		}
		if detail.BidWinnerName != nil {
			t.Fatalf("expected nil winner name, got %q", *detail.BidWinnerName)
		}
		if detail.Auction.BidWinnerID != nil {
			t.Fatalf("expected nil winner id, got %d", *detail.Auction.BidWinnerID)
		}
		if detail.TotalBids != 0 || detail.TotalBidders != 0 {
			t.Fatalf("expected no bids, got %d bids from %d bidders", detail.TotalBids, detail.TotalBidders)
		}
	})

	t.Run("auction not found", func(t *testing.T) {
		detail, err := repo.FindByID(ctx, 999999)
		if detail != nil {
			t.Fatalf("expected nil detail, got %+v", detail)
		}
		if !errors.Is(err, ErrAuctionNotFound) {
			t.Fatalf("expected error %v, got %v", ErrAuctionNotFound, err)
		}
	})

	t.Run("top bids are ordered by price then by time", func(t *testing.T) {
		topBids, err := repo.FindTopBids(ctx, running, 3)
		if err != nil {
			t.Fatalf("find top bids: %v", err)
		}
		if len(topBids) != 3 {
			t.Fatalf("expected the limit to be respected, got %d bids", len(topBids))
		}
		if !topBids[0].Bid.BidPrice.Equal(decimal.NewFromInt(2000000)) || !topBids[1].Bid.BidPrice.Equal(decimal.NewFromInt(2000000)) {
			t.Fatalf("expected the two highest bids first, got %s and %s", topBids[0].Bid.BidPrice, topBids[1].Bid.BidPrice)
		}
		// same price: the earlier bid wins
		if topBids[0].UserName != "Siti Rahma" || topBids[1].UserName != "John Doe" {
			t.Fatalf("expected the earlier bid to rank first on a tie, got %q then %q", topBids[0].UserName, topBids[1].UserName)
		}
		if !topBids[2].Bid.BidPrice.Equal(decimal.NewFromInt(1500000)) {
			t.Fatalf("expected third bid to be 1500000, got %s", topBids[2].Bid.BidPrice)
		}
	})

	t.Run("top bids for an auction without bids", func(t *testing.T) {
		topBids, err := repo.FindTopBids(ctx, expired, 3)
		if err != nil {
			t.Fatalf("find top bids: %v", err)
		}
		if len(topBids) != 0 {
			t.Fatalf("expected no bids, got %d", len(topBids))
		}
	})
}
