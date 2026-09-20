package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
	svc "github.com/grup-baru-belajar/auction-bid-repo/internal/services"
	"github.com/shopspring/decimal"
)

type mockAuctionDetailRepo struct {
	detail  *repository.AuctionDetail
	topBids []repository.TopBid
	findErr error
	topErr  error

	gotID       int64
	gotTopLimit int
}

func (m *mockAuctionDetailRepo) FindByID(ctx context.Context, id int64) (*repository.AuctionDetail, error) {
	m.gotID = id
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.detail, nil
}

func (m *mockAuctionDetailRepo) FindTopBids(ctx context.Context, auctionID int64, limit int) ([]repository.TopBid, error) {
	m.gotTopLimit = limit
	if m.topErr != nil {
		return nil, m.topErr
	}
	return m.topBids, nil
}

// newAuctionDetail builds the row shape the repository would return
func newAuctionDetail(winnerID *int64, winnerName *string) *repository.AuctionDetail {
	return &repository.AuctionDetail{
		Auction: models.Auction{
			ID:            1,
			AuctionName:   "iPhone 15 Pro",
			Description:   "Brand new",
			ImageLink:     "https://example.com/a.jpg",
			BidWinnerID:   winnerID,
			StartingPrice: decimal.NewFromInt(10000000),
			LastPrice:     decimal.NewFromInt(16000000),
			CreatedAt:     time.Now().Add(-48 * time.Hour),
			EndTime:       time.Now().Add(48 * time.Hour),
			IsCompleted:   false,
		},
		BidWinnerName: winnerName,
		TotalBids:     6,
		TotalBidders:  5,
	}
}

func TestGetAuctionByID_Success(t *testing.T) {
	winnerID := int64(5)
	winnerName := "Siti Rahma"
	mock := &mockAuctionDetailRepo{
		detail: newAuctionDetail(&winnerID, &winnerName),
		topBids: []repository.TopBid{
			{Bid: models.Bid{ID: 9, AuctionID: 1, UserID: 5, BidPrice: decimal.NewFromInt(16000000)}, UserName: "Siti Rahma"},
			{Bid: models.Bid{ID: 8, AuctionID: 1, UserID: 2, BidPrice: decimal.NewFromInt(15000000)}, UserName: "John Doe"},
		},
	}
	s := svc.NewAuctionDetailService(mock)

	got, err := s.GetAuctionByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.gotID != 1 {
		t.Fatalf("expected repo called with id 1, got %d", mock.gotID)
	}
	if got.ID != 1 || got.AuctionName != "iPhone 15 Pro" || got.TotalBids != 6 || got.TotalBidders != 5 {
		t.Fatalf("unexpected auction detail response: %+v", got)
	}
	if got.BidWinner == nil || got.BidWinner.ID != winnerID || got.BidWinner.Name != winnerName {
		t.Fatalf("unexpected bid winner: %+v", got.BidWinner)
	}
	if len(got.TopBids) != 2 || got.TopBids[0].ID != 9 || got.TopBids[0].UserName != "Siti Rahma" {
		t.Fatalf("unexpected top bids: %+v", got.TopBids)
	}
}

func TestGetAuctionByID_TopBidsLimit(t *testing.T) {
	mock := &mockAuctionDetailRepo{detail: newAuctionDetail(nil, nil)}
	s := svc.NewAuctionDetailService(mock)

	if _, err := s.GetAuctionByID(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.gotTopLimit != 3 {
		t.Fatalf("expected detail view to ask for 3 top bids, got %d", mock.gotTopLimit)
	}
}

func TestGetAuctionByID_NoWinner(t *testing.T) {
	tests := []struct {
		name       string
		winnerID   *int64
		winnerName *string
	}{
		{"no winner at all", nil, nil},
		{"winner id without name", func() *int64 { id := int64(5); return &id }(), nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockAuctionDetailRepo{detail: newAuctionDetail(tc.winnerID, tc.winnerName)}
			s := svc.NewAuctionDetailService(mock)

			got, err := s.GetAuctionByID(context.Background(), 1)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.BidWinner != nil {
				t.Fatalf("expected nil bid winner, got %+v", got.BidWinner)
			}
			if got.TopBids == nil {
				t.Fatalf("expected empty slice so json renders [] instead of null")
			}
		})
	}
}

func TestGetAuctionByID_Errors(t *testing.T) {
	tests := []struct {
		name    string
		findErr error
		topErr  error
		wantErr error
	}{
		{"auction not found", repository.ErrAuctionNotFound, nil, svc.ErrAuctionNotFound},
		{"find by id fails", errors.New("db fail"), nil, nil},
		{"find top bids fails", nil, errors.New("db fail"), nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockAuctionDetailRepo{
				detail:  newAuctionDetail(nil, nil),
				findErr: tc.findErr,
				topErr:  tc.topErr,
			}
			s := svc.NewAuctionDetailService(mock)

			got, err := s.GetAuctionByID(context.Background(), 1)
			if got != nil {
				t.Fatalf("expected nil response when error, got %+v", got)
			}
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("expected error %v, got %v", tc.wantErr, err)
				}
			} else {
				// expect wrapped error
				if err == nil {
					t.Fatalf("expected non-nil error for repo error")
				}
				if errors.Is(err, svc.ErrAuctionNotFound) {
					t.Fatalf("db failure must not be reported as auction not found")
				}
			}
		})
	}
}
