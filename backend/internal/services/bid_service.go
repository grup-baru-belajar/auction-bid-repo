package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
	"github.com/shopspring/decimal"
)

var (
	ErrAuctionCompleted = errors.New("auction already completed")
	ErrBidTooLow        = errors.New("bid price must be greater than current price")
)

type bidRepository interface {
	Place(ctx context.Context, auctionID, userID int64, bidPrice decimal.Decimal) (*models.Bid, error)
}

type BidService interface {
	PlaceBid(ctx context.Context, userID int64, req models.BidRequest) (*models.BidResponse, error)
}

type bidService struct {
	repo bidRepository
}

func NewBidService(repo bidRepository) BidService {
	return &bidService{repo: repo}
}

func (s *bidService) PlaceBid(ctx context.Context, userID int64, req models.BidRequest) (*models.BidResponse, error) {
	bid, err := s.repo.Place(ctx, req.AuctionID, userID, req.BidPrice)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAuctionNotFound):
			return nil, ErrAuctionNotFound
		case errors.Is(err, repository.ErrAuctionCompleted):
			return nil, ErrAuctionCompleted
		case errors.Is(err, repository.ErrBidTooLow):
			return nil, ErrBidTooLow
		default:
			return nil, fmt.Errorf("place bid service: %w", err)
		}
	}

	return &models.BidResponse{
		ID:        bid.ID,
		AuctionID: bid.AuctionID,
		UserID:    bid.UserID,
		BidPrice:  bid.BidPrice,
		CreatedAt: bid.CreatedAt,
	}, nil
}
