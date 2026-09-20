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

type mockBidRepo struct {
    result *models.Bid
    err    error

    gotAuctionID int64
    gotUserID    int64
    gotBidPrice  decimal.Decimal
}

func (m *mockBidRepo) Place(ctx context.Context, auctionID, userID int64, bidPrice decimal.Decimal) (*models.Bid, error) {
    m.gotAuctionID = auctionID
    m.gotUserID = userID
    m.gotBidPrice = bidPrice
    if m.err != nil {
        return nil, m.err
    }
    return m.result, nil
}

func TestPlaceBid_Success(t *testing.T) {
    now := time.Now()
    bid := &models.Bid{
        ID:        10,
        AuctionID: 1,
        UserID:    2,
        BidPrice:  decimal.NewFromInt(15000000),
        CreatedAt: now,
    }

    mock := &mockBidRepo{result: bid}
    s := svc.NewBidService(mock)

    req := models.BidRequest{AuctionID: 1, BidPrice: decimal.NewFromInt(15000000)}
    got, err := s.PlaceBid(context.Background(), 2, req)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if got.ID != bid.ID || got.AuctionID != bid.AuctionID || got.UserID != bid.UserID {
        t.Fatalf("unexpected bid response: %+v", got)
    }
}

func TestPlaceBid_Errors(t *testing.T) {
    tests := []struct{
        name string
        repoErr error
        wantErr error
    }{
        {"auction not found", repository.ErrAuctionNotFound, svc.ErrAuctionNotFound},
        {"auction completed", repository.ErrAuctionCompleted, svc.ErrAuctionCompleted},
        {"bid too low", repository.ErrBidTooLow, svc.ErrBidTooLow},
        {"other repo error", errors.New("db fail"), nil},
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            mock := &mockBidRepo{err: tc.repoErr}
            s := svc.NewBidService(mock)
            req := models.BidRequest{AuctionID: 1, BidPrice: decimal.NewFromInt(100)}
            got, err := s.PlaceBid(context.Background(), 1, req)
            if tc.wantErr != nil {
                if !errors.Is(err, tc.wantErr) {
                    t.Fatalf("expected error %v, got %v", tc.wantErr, err)
                }
                if got != nil {
                    t.Fatalf("expected nil response when error, got %+v", got)
                }
            } else {
                // expect wrapped error
                if err == nil {
                    t.Fatalf("expected non-nil error for repo error")
                }
            }
        })
    }
}
