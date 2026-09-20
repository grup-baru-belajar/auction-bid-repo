package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/middlewares"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
	svc "github.com/grup-baru-belajar/auction-bid-repo/internal/services"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/token"
	wsh "github.com/grup-baru-belajar/auction-bid-repo/internal/websocket"
	"github.com/shopspring/decimal"
)

type mockBidService struct {
	resp *models.BidResponse
	err  error
}

func (m *mockBidService) PlaceBid(ctx context.Context, userID int64, req models.BidRequest) (*models.BidResponse, error) {
	return m.resp, m.err
}

type mockAuctionDetailRepo struct{}

func (m *mockAuctionDetailRepo) FindByID(ctx context.Context, id int64) (*repository.AuctionDetail, error) {
	// return minimal valid detail for websocket buildMessage
	name := "winner"
	return &repository.AuctionDetail{
		Auction:       repository.AuctionDetail{}.Auction,
		BidWinnerName: &name,
		TotalBids:     1,
		TotalBidders:  1,
	}, nil
}

func (m *mockAuctionDetailRepo) FindTopBids(ctx context.Context, auctionID int64, limit int) ([]repository.TopBid, error) {
	return []repository.TopBid{}, nil
}

func TestPostBid_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// prepare mock bid service
	now := time.Now()
	resp := &models.BidResponse{ID: 5, AuctionID: 1, UserID: 2, BidPrice: decimal.NewFromInt(15000000), CreatedAt: now}
	mockSvc := &mockBidService{resp: resp}

	// prepare ws handler (needs hub + repo)
	hub := wsh.NewHub()
	repo := &mockAuctionDetailRepo{}
	wsHandler := wsh.NewHandler(hub, repo)

	h := handlers.New(nil, nil, nil, mockSvc, nil, wsHandler, nil)

	body := `{"auctionId":1,"bidPrice":"15000000"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bid", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// inject claims as middleware would do
	c.Set(middlewares.ClaimsKey, &token.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "2"}})

	h.PostBid(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusCreated, w.Code, w.Body.String())
	}

	var ar models.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &ar); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !ar.Success {
		t.Fatalf("expected success true, got false, body=%s", w.Body.String())
	}
}

func TestPostBid_Handler_ErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		svcErr   error
		wantCode int
	}{
		{"auction not found", svc.ErrAuctionNotFound, http.StatusNotFound},
		{"auction completed", svc.ErrAuctionCompleted, http.StatusConflict},
		{"bid too low", svc.ErrBidTooLow, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockBidService{resp: nil, err: tc.svcErr}
			hub := wsh.NewHub()
			repo := &mockAuctionDetailRepo{}
			wsHandler := wsh.NewHandler(hub, repo)
			h := handlers.New(nil, nil, nil, mockSvc, nil, wsHandler, nil)

			body := `{"auctionId":1,"bidPrice":"100"}`
			req := httptest.NewRequest(http.MethodPost, "/api/v1/bid", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Set(middlewares.ClaimsKey, &token.Claims{RegisteredClaims: jwt.RegisteredClaims{Subject: "1"}})

			h.PostBid(c)

			if w.Code != tc.wantCode {
				t.Fatalf("expected status %d, got %d, body=%s", tc.wantCode, w.Code, w.Body.String())
			}
		})
	}
}
