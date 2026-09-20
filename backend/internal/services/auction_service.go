package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
)

var (
	ErrInvalidStartingPrice = errors.New("starting price must be greater than zero")
	ErrInvalidEndTime       = errors.New("end time must be in the future")
)

type auctionRepository interface {
	Create(ctx context.Context, auction *models.Auction) (*models.Auction, error)
	FindAll(ctx context.Context, limit, offset int, isCompleted *bool, search string) ([]repository.AuctionWithWinner, int64, error)
}

type AuctionService interface {
	CreateAuction(ctx context.Context, req models.AuctionRequest) (*models.AuctionResponse, error)
	GetAuctions(ctx context.Context, query models.GetAuctionsQuery) ([]models.AuctionResponse, models.PaginationResponse, error)
}

type auctionService struct {
	repo auctionRepository
}

func NewAuctionService(repo auctionRepository) AuctionService {
	return &auctionService{repo: repo}
}

func (s *auctionService) CreateAuction(ctx context.Context, req models.AuctionRequest) (*models.AuctionResponse, error) {
	if !req.StartingPrice.IsPositive() {
		return nil, ErrInvalidStartingPrice
	}

	if !req.EndTime.After(time.Now()) {
		return nil, ErrInvalidEndTime
	}

	auction := &models.Auction{
		AuctionName:   req.AuctionName,
		Description:   req.Description,
		ImageLink:     req.ImageLink,
		StartingPrice: req.StartingPrice,
		LastPrice:     req.StartingPrice,
		EndTime:       req.EndTime,
		IsCompleted:   false,
	}

	created, err := s.repo.Create(ctx, auction)
	if err != nil {
		return nil, fmt.Errorf("create auction service: %w", err)
	}

	return &models.AuctionResponse{
		ID:            created.ID,
		AuctionName:   created.AuctionName,
		Description:   created.Description,
		ImageLink:     created.ImageLink,
		StartingPrice: created.StartingPrice,
		LastPrice:     created.LastPrice,
		CreatedAt:     created.CreatedAt,
		EndTime:       created.EndTime,
		IsCompleted:   created.IsCompleted,
		BidWinner:     nil,
	}, nil
}

func (s *auctionService) GetAuctions(ctx context.Context, query models.GetAuctionsQuery) ([]models.AuctionResponse, models.PaginationResponse, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}

	limit := query.Limit
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	items, total, err := s.repo.FindAll(ctx, limit, offset, query.IsCompleted, query.Search)
	if err != nil {
		return nil, models.PaginationResponse{}, fmt.Errorf("get auctions service: %w", err)
	}

	totalPages := int64(math.Ceil(float64(total) / float64(limit)))

	responses := make([]models.AuctionResponse, 0, len(items))
	for _, item := range items {
		var bidWinner *models.BidWinnerResponse
		if item.Auction.BidWinnerID != nil && item.BidWinnerName != nil {
			bidWinner = &models.BidWinnerResponse{
				ID:   *item.Auction.BidWinnerID,
				Name: *item.BidWinnerName,
			}
		}

		responses = append(responses, models.AuctionResponse{
			ID:            item.Auction.ID,
			AuctionName:   item.Auction.AuctionName,
			Description:   item.Auction.Description,
			ImageLink:     item.Auction.ImageLink,
			StartingPrice: item.Auction.StartingPrice,
			LastPrice:     item.Auction.LastPrice,
			CreatedAt:     item.Auction.CreatedAt,
			EndTime:       item.Auction.EndTime,
			IsCompleted:   item.Auction.IsCompleted,
			BidWinner:     bidWinner,
		})
	}

	pagination := models.PaginationResponse{
		Page:       int64(page),
		Limit:      int64(limit),
		Total:      total,
		TotalPages: totalPages,
	}

	return responses, pagination, nil
}
