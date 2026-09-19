package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
)

var ErrAuctionNotFound = errors.New("auction not found")

const topBidsLimit = 3

type AuctionDetailService interface {
	GetAuctionByID(ctx context.Context, id int64) (*models.AuctionDetailResponse, error)
}

type auctionDetailService struct {
	repo repository.AuctionDetailRepository
}

func NewAuctionDetailService(repo repository.AuctionDetailRepository) AuctionDetailService {
	return &auctionDetailService{repo: repo}
}

func (s *auctionDetailService) GetAuctionByID(ctx context.Context, id int64) (*models.AuctionDetailResponse, error) {
	detail, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrAuctionNotFound) {
			return nil, ErrAuctionNotFound
		}
		return nil, fmt.Errorf("get auction by id service: %w", err)
	}

	topBids, err := s.repo.FindTopBids(ctx, id, topBidsLimit)
	if err != nil {
		return nil, fmt.Errorf("get top bids service: %w", err)
	}

	var bidWinner *models.BidWinnerResponse
	if detail.Auction.BidWinnerID != nil && detail.BidWinnerName != nil {
		bidWinner = &models.BidWinnerResponse{
			ID:   *detail.Auction.BidWinnerID,
			Name: *detail.BidWinnerName,
		}
	}

	topBidResponses := make([]models.TopBidResponse, 0, len(topBids))
	for _, tb := range topBids {
		topBidResponses = append(topBidResponses, models.TopBidResponse{
			ID:        tb.Bid.ID,
			UserID:    tb.Bid.UserID,
			UserName:  tb.UserName,
			BidPrice:  tb.Bid.BidPrice,
			CreatedAt: tb.Bid.CreatedAt,
		})
	}

	return &models.AuctionDetailResponse{
		ID:            detail.Auction.ID,
		AuctionName:   detail.Auction.AuctionName,
		Description:   detail.Auction.Description,
		ImageLink:     detail.Auction.ImageLink,
		StartingPrice: detail.Auction.StartingPrice,
		LastPrice:     detail.Auction.LastPrice,
		CreatedAt:     detail.Auction.CreatedAt,
		EndTime:       detail.Auction.EndTime,
		IsCompleted:   detail.Auction.IsCompleted,
		BidWinner:     bidWinner,
		TotalBids:     detail.TotalBids,
		TotalBidders:  detail.TotalBidders,
		TopBids:       topBidResponses,
	}, nil
}