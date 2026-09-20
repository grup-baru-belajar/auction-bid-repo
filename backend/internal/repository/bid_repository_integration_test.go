//go:build integration
// +build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

func TestBidRepository_Place_EndToEnd(t *testing.T) {
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

	// create minimal tables required by bid repository
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
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        end_time TIMESTAMP NOT NULL,
        is_completed BOOLEAN DEFAULT FALSE
    );

    CREATE TABLE bids (
        id BIGSERIAL PRIMARY KEY,
        auction_id BIGINT NOT NULL,
        user_id BIGINT NOT NULL,
        bid_price NUMERIC(18,2) NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT now()
    );
    `

	if _, err := db.ExecContext(context.Background(), createSQL); err != nil {
		t.Fatalf("create tables: %v", err)
	}

	// insert test user and auction
	var userID int64
	err = db.QueryRowContext(context.Background(), `INSERT INTO users (name) VALUES ($1) RETURNING id`, "testuser").Scan(&userID)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}

	endTime := time.Now().Add(1 * time.Hour)
	var auctionID int64
	err = db.QueryRowContext(context.Background(), `INSERT INTO auctions (auction_name, description, image_link, starting_price, last_price, end_time) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`, "testauction", "desc", "img", decimal.NewFromInt(1000000), decimal.NewFromInt(1000000), endTime).Scan(&auctionID)
	if err != nil {
		t.Fatalf("insert auction: %v", err)
	}

	repo := NewBidRepository(db)

	// place a higher bid
	bidPrice := decimal.NewFromInt(1500000)
	bid, err := repo.Place(context.Background(), auctionID, userID, bidPrice)
	if err != nil {
		t.Fatalf("place bid: %v", err)
	}
	if bid.AuctionID != auctionID || bid.UserID != userID {
		t.Fatalf("unexpected bid returned: %+v", bid)
	}

	// verify auction last_price and bid_winner_id updated
	var lastPrice decimal.Decimal
	var bidWinnerID sql.NullInt64
	err = db.QueryRowContext(context.Background(), `SELECT last_price, bid_winner_id FROM auctions WHERE id = $1`, auctionID).Scan(&lastPrice, &bidWinnerID)
	if err != nil {
		t.Fatalf("select auction after bid: %v", err)
	}
	if !lastPrice.Equal(bidPrice) {
		t.Fatalf("expected auction.last_price %s got %s", bidPrice.String(), lastPrice.String())
	}
	if !bidWinnerID.Valid || bidWinnerID.Int64 != userID {
		t.Fatalf("expected bid_winner_id %d got %v", userID, bidWinnerID)
	}
}
