package services

import (
	"context"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
)

type ReportingService interface {
	GetTopAuction(ctx context.Context, limit int) ([]models.TopAuction, error)
	GetAuctionActivity(ctx context.Context, intervalDays int) ([]models.AuctionActivity, error)
	GetAuctionStatus(ctx context.Context) ([]models.AuctionStatus, error)
	GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error)
	GetTotalTransaction(ctx context.Context, interval string) (models.TotalTransaction, error)
	GetAuctionSummary(ctx context.Context) (models.AuctionSummary, error)
	GetTransactionOverview(ctx context.Context, weeks int) ([]models.TransactionWeek, error)
	GetTopBidders(ctx context.Context, limit int, sortBy string) ([]models.TopBidder, error)
}

type reportingService struct {
	repo repository.ReportingRepository
}

func NewReportingService(repo repository.ReportingRepository) ReportingService {
	return &reportingService{repo: repo}
}

func (s *reportingService) GetTopAuction(ctx context.Context, limit int) ([]models.TopAuction, error) {
	return s.repo.GetTopAuction(ctx, limit)
}

func (s *reportingService) GetAuctionActivity(ctx context.Context, intervalDays int) ([]models.AuctionActivity, error) {
	return s.repo.GetAuctionActivity(ctx, intervalDays)
}

func (s *reportingService) GetTransactionOverview(ctx context.Context, weeks int) ([]models.TransactionWeek, error) {
	return s.repo.GetTransactionOverview(ctx, weeks)
}

func (s *reportingService) GetAuctionStatus(ctx context.Context) ([]models.AuctionStatus, error) {
	return s.repo.GetAuctionStatus(ctx)
}

func (s *reportingService) GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error) {
	return s.repo.GetTotalBidders(ctx, interval, auctionId)
}

func (s *reportingService) GetTotalTransaction(ctx context.Context, interval string) (models.TotalTransaction, error) {
	return s.repo.GetTotalTransaction(ctx, interval)
}
func (s *reportingService) GetAuctionSummary(ctx context.Context) (models.AuctionSummary, error) {
	return s.repo.GetAuctionSummary(ctx)
}

func (s *reportingService) GetTopBidders(ctx context.Context, limit int, sortBy string) ([]models.TopBidder, error) {
	return s.repo.GetTopBidders(ctx, limit, sortBy)
}
