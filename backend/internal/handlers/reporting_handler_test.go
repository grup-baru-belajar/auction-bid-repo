package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

// mockReportingService records the args it was called with so tests can
// assert how the handler parses query params, and returns canned results/errors.
type mockReportingService struct {
	topAuction      []models.TopAuction
	auctionActivity []models.AuctionActivity
	auctionStatus   []models.AuctionStatus
	totalBidders    int
	totalTx         models.TotalTransaction
	summary         models.AuctionSummary
	txOverview      []models.TransactionWeek
	topBidders      []models.TopBidder
	err             error

	gotTopAuctionLimit      int
	gotActivityIntervalDays int
	gotBiddersInterval      string
	gotBiddersAuctionID     string
	gotTxInterval           string
	gotOverviewWeeks        int
	gotTopBiddersLimit      int
	gotTopBiddersSortBy     string
}

func (m *mockReportingService) GetTopAuction(ctx context.Context, limit int) ([]models.TopAuction, error) {
	m.gotTopAuctionLimit = limit
	if m.err != nil {
		return nil, m.err
	}
	return m.topAuction, nil
}

func (m *mockReportingService) GetAuctionActivity(ctx context.Context, intervalDays int) ([]models.AuctionActivity, error) {
	m.gotActivityIntervalDays = intervalDays
	if m.err != nil {
		return nil, m.err
	}
	return m.auctionActivity, nil
}

func (m *mockReportingService) GetAuctionStatus(ctx context.Context) ([]models.AuctionStatus, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.auctionStatus, nil
}

func (m *mockReportingService) GetTotalBidders(ctx context.Context, interval string, auctionId string) (int, error) {
	m.gotBiddersInterval = interval
	m.gotBiddersAuctionID = auctionId
	if m.err != nil {
		return 0, m.err
	}
	return m.totalBidders, nil
}

func (m *mockReportingService) GetTotalTransaction(ctx context.Context, interval string) (models.TotalTransaction, error) {
	m.gotTxInterval = interval
	if m.err != nil {
		return models.TotalTransaction{}, m.err
	}
	return m.totalTx, nil
}

func (m *mockReportingService) GetAuctionSummary(ctx context.Context) (models.AuctionSummary, error) {
	if m.err != nil {
		return models.AuctionSummary{}, m.err
	}
	return m.summary, nil
}

func (m *mockReportingService) GetTransactionOverview(ctx context.Context, weeks int) ([]models.TransactionWeek, error) {
	m.gotOverviewWeeks = weeks
	if m.err != nil {
		return nil, m.err
	}
	return m.txOverview, nil
}

func (m *mockReportingService) GetTopBidders(ctx context.Context, limit int, sortBy string) ([]models.TopBidder, error) {
	m.gotTopBiddersLimit = limit
	m.gotTopBiddersSortBy = sortBy
	if m.err != nil {
		return nil, m.err
	}
	return m.topBidders, nil
}

// newReportingTestContext builds a gin context for the given method/path,
// wired up so tests can inspect status code and response body.
func newReportingTestContext(path string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func decodeAPIResponse(t *testing.T, w *httptest.ResponseRecorder) models.APIResponse {
	t.Helper()
	var ar models.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &ar); err != nil {
		t.Fatalf("failed to unmarshal response: %v, body=%s", err, w.Body.String())
	}
	return ar
}

// --- GET /reporting/top-auction ---

func TestGetTopAuction_Handler_Success(t *testing.T) {
	mockSvc := &mockReportingService{topAuction: []models.TopAuction{
		{ID: 1, AuctionName: "iPhone 15 Pro", TotalBids: 12, Bidders: 5, Status: "ACTIVE"},
	}}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/top-auction")
	h.GetTopAuction(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotTopAuctionLimit != 5 {
		t.Fatalf("expected service called with limit 5, got %d", mockSvc.gotTopAuctionLimit)
	}
	ar := decodeAPIResponse(t, w)
	if !ar.Success {
		t.Fatalf("expected success true, body=%s", w.Body.String())
	}
}

func TestGetTopAuction_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("pq: connection refused")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/top-auction")
	h.GetTopAuction(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
	ar := decodeAPIResponse(t, w)
	if ar.Success {
		t.Fatalf("expected success false, body=%s", w.Body.String())
	}
}

// --- GET /reporting/auction-activity ---

func TestGetAuctionActivity_Handler_DefaultInterval(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-activity")
	h.GetAuctionActivity(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotActivityIntervalDays != 7 {
		t.Fatalf("expected default interval 7, got %d", mockSvc.gotActivityIntervalDays)
	}
}

func TestGetAuctionActivity_Handler_CustomInterval(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-activity?interval=30")
	h.GetAuctionActivity(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotActivityIntervalDays != 30 {
		t.Fatalf("expected interval 30, got %d", mockSvc.gotActivityIntervalDays)
	}
}

func TestGetAuctionActivity_Handler_InvalidIntervalFallsBackToDefault(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{"not a number", "interval=abc"},
		{"zero", "interval=0"},
		{"negative", "interval=-5"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockReportingService{}
			h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

			c, w := newReportingTestContext("/api/v1/reporting/auction-activity?" + tc.query)
			h.GetAuctionActivity(c)

			if w.Code != http.StatusOK {
				t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
			}
			if mockSvc.gotActivityIntervalDays != 7 {
				t.Fatalf("expected fallback to default interval 7, got %d", mockSvc.gotActivityIntervalDays)
			}
		})
	}
}

func TestGetAuctionActivity_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("db fail")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-activity")
	h.GetAuctionActivity(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

// --- GET /reporting/auction-status ---

func TestGetAuctionStatus_Handler_Success(t *testing.T) {
	mockSvc := &mockReportingService{auctionStatus: []models.AuctionStatus{
		{Status: "ACTIVE", Total: 3}, {Status: "ENDED", Total: 7},
	}}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-status")
	h.GetAuctionStatus(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	ar := decodeAPIResponse(t, w)
	data, ok := ar.Data.([]any)
	if !ok || len(data) != 2 {
		t.Fatalf("expected 2 status rows, body=%s", w.Body.String())
	}
}

func TestGetAuctionStatus_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("db fail")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-status")
	h.GetAuctionStatus(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

// --- GET /reporting/total-bidders ---

func TestGetTotalBidders_Handler_DefaultInterval(t *testing.T) {
	mockSvc := &mockReportingService{totalBidders: 42}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-bidders")
	h.GetTotalBidders(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotBiddersInterval != "7 days" {
		t.Fatalf("expected default interval %q, got %q", "7 days", mockSvc.gotBiddersInterval)
	}
}

func TestGetTotalBidders_Handler_AllInterval(t *testing.T) {
	mockSvc := &mockReportingService{totalBidders: 100}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-bidders?interval=all")
	h.GetTotalBidders(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotBiddersInterval != "all" {
		t.Fatalf("expected interval %q to pass through unchanged, got %q", "all", mockSvc.gotBiddersInterval)
	}
}

func TestGetTotalBidders_Handler_CustomIntervalAndAuctionID(t *testing.T) {
	mockSvc := &mockReportingService{totalBidders: 3}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-bidders?interval=30&auctionId=7")
	h.GetTotalBidders(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotBiddersInterval != "30 days" {
		t.Fatalf("expected interval %q, got %q", "30 days", mockSvc.gotBiddersInterval)
	}
	if mockSvc.gotBiddersAuctionID != "7" {
		t.Fatalf("expected auctionId %q, got %q", "7", mockSvc.gotBiddersAuctionID)
	}
}

func TestGetTotalBidders_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("db fail")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-bidders")
	h.GetTotalBidders(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

// --- GET /reporting/total-transaction ---

func TestGetTotalTransaction_Handler_DefaultInterval(t *testing.T) {
	mockSvc := &mockReportingService{totalTx: models.TotalTransaction{TransactionCount: "10"}}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-transaction")
	h.GetTotalTransaction(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotTxInterval != "7 days" {
		t.Fatalf("expected default interval %q, got %q", "7 days", mockSvc.gotTxInterval)
	}
}

func TestGetTotalTransaction_Handler_CustomInterval(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-transaction?interval=14")
	h.GetTotalTransaction(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotTxInterval != "14 days" {
		t.Fatalf("expected interval %q, got %q", "14 days", mockSvc.gotTxInterval)
	}
}

func TestGetTotalTransaction_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("db fail")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/total-transaction")
	h.GetTotalTransaction(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

// --- GET /reporting/auction-summary ---

func TestGetAuctionSummary_Handler_Success(t *testing.T) {
	mockSvc := &mockReportingService{summary: models.AuctionSummary{
		TotalAuctions: 10, OngoingAuctions: 4, CompletedAuctions: 6,
		TotalBidsOngoing: 20, TotalBidsCompleted: 80, TotalBidsAll: 100,
	}}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-summary")
	h.GetAuctionSummary(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	ar := decodeAPIResponse(t, w)
	data, ok := ar.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected summary object, body=%s", w.Body.String())
	}
	if data["totalAuctions"] != float64(10) || data["totalBidsAll"] != float64(100) {
		t.Fatalf("unexpected summary payload: %+v", data)
	}
}

func TestGetAuctionSummary_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("db fail")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/auction-summary")
	h.GetAuctionSummary(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

// --- GET /reporting/transaction-overview ---

func TestGetTransactionOverview_Handler_DefaultWeeks(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/transaction-overview")
	h.GetTransactionOverview(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotOverviewWeeks != 4 {
		t.Fatalf("expected default weeks 4, got %d", mockSvc.gotOverviewWeeks)
	}
}

func TestGetTransactionOverview_Handler_CustomWeeks(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/transaction-overview?weeks=8")
	h.GetTransactionOverview(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotOverviewWeeks != 8 {
		t.Fatalf("expected weeks 8, got %d", mockSvc.gotOverviewWeeks)
	}
}

func TestGetTransactionOverview_Handler_InvalidWeeksFallsBackToDefault(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/transaction-overview?weeks=-1")
	h.GetTransactionOverview(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotOverviewWeeks != 4 {
		t.Fatalf("expected fallback to default weeks 4, got %d", mockSvc.gotOverviewWeeks)
	}
}

func TestGetTransactionOverview_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("db fail")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/transaction-overview")
	h.GetTransactionOverview(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}

// --- GET /reporting/top-bidders ---

func TestGetTopBidders_Handler_Defaults(t *testing.T) {
	mockSvc := &mockReportingService{topBidders: []models.TopBidder{{UserID: 1, Username: "siti"}}}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/top-bidders")
	h.GetTopBidders(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotTopBiddersLimit != 5 {
		t.Fatalf("expected default limit 5, got %d", mockSvc.gotTopBiddersLimit)
	}
	if mockSvc.gotTopBiddersSortBy != "amount" {
		t.Fatalf("expected default sort_by %q, got %q", "amount", mockSvc.gotTopBiddersSortBy)
	}
}

func TestGetTopBidders_Handler_CustomParams(t *testing.T) {
	mockSvc := &mockReportingService{topBidders: []models.TopBidder{}}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/top-bidders?limit=10&sort_by=transaction_count")
	h.GetTopBidders(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotTopBiddersLimit != 10 {
		t.Fatalf("expected limit 10, got %d", mockSvc.gotTopBiddersLimit)
	}
	if mockSvc.gotTopBiddersSortBy != "transaction_count" {
		t.Fatalf("expected sort_by %q, got %q", "transaction_count", mockSvc.gotTopBiddersSortBy)
	}
}

func TestGetTopBidders_Handler_InvalidLimit(t *testing.T) {
	mockSvc := &mockReportingService{}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/top-bidders?limit=abc")
	h.GetTopBidders(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}
	ar := decodeAPIResponse(t, w)
	if ar.Message != "Invalid limit parameter" {
		t.Fatalf("expected message %q, got %q", "Invalid limit parameter", ar.Message)
	}
}

func TestGetTopBidders_Handler_ServiceError(t *testing.T) {
	mockSvc := &mockReportingService{err: errors.New("invalid sortBy value: bogus")}
	h := handlers.New(nil, nil, nil, nil, mockSvc, nil, nil)

	c, w := newReportingTestContext("/api/v1/reporting/top-bidders?sort_by=bogus")
	h.GetTopBidders(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusInternalServerError, w.Code, w.Body.String())
	}
}
