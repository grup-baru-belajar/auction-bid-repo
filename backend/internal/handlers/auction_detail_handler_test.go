package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/handlers"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	svc "github.com/grup-baru-belajar/auction-bid-repo/internal/services"
)

type mockAuctionDetailService struct {
	resp *models.AuctionDetailResponse
	err  error

	called bool
	gotID  int64
}

func (m *mockAuctionDetailService) GetAuctionByID(ctx context.Context, id int64) (*models.AuctionDetailResponse, error) {
	m.called = true
	m.gotID = id
	return m.resp, m.err
}

func TestGetAuctionByID_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// prepare mock auction detail service
	resp := &models.AuctionDetailResponse{ID: 1, AuctionName: "iPhone 15 Pro", TotalBids: 6, TotalBidders: 5}
	mockSvc := &mockAuctionDetailService{resp: resp}

	h := handlers.New(nil, nil, mockSvc, nil, nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auctions/1", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	h.GetAuctionByID(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
	}
	if mockSvc.gotID != 1 {
		t.Fatalf("expected service called with id 1, got %d", mockSvc.gotID)
	}

	var ar models.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &ar); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if !ar.Success {
		t.Fatalf("expected success true, got false, body=%s", w.Body.String())
	}

	data, ok := ar.Data.(map[string]any)
	if !ok {
		t.Fatalf("expected auction detail object, body=%s", w.Body.String())
	}
	if data["auctionName"] != "iPhone 15 Pro" {
		t.Fatalf("unexpected auction payload: %+v", data)
	}
}

func TestGetAuctionByID_Handler_ErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		svcErr   error
		wantCode int
	}{
		{"auction not found", svc.ErrAuctionNotFound, http.StatusNotFound},
		{"unexpected error", errors.New("pq: connection refused on 10.0.0.5"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAuctionDetailService{resp: nil, err: tc.svcErr}
			h := handlers.New(nil, nil, mockSvc, nil, nil, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/auctions/1", nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			h.GetAuctionByID(c)

			if w.Code != tc.wantCode {
				t.Fatalf("expected status %d, got %d, body=%s", tc.wantCode, w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "10.0.0.5") {
				t.Fatalf("internal error detail must not leak to the client, body=%s", w.Body.String())
			}
		})
	}
}

func TestGetAuctionByID_Handler_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name string
		id   string
	}{
		{"not a number", "abc"},
		{"zero", "0"},
		{"negative", "-1"},
		{"decimal", "1.5"},
		{"empty", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAuctionDetailService{}
			h := handlers.New(nil, nil, mockSvc, nil, nil, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/auctions/"+tc.id, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tc.id}}

			h.GetAuctionByID(c)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
			}

			var ar models.APIResponse
			if err := json.Unmarshal(w.Body.Bytes(), &ar); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if ar.Message != "Invalid auction id" {
				t.Fatalf("expected message %q, got %q", "Invalid auction id", ar.Message)
			}
			if mockSvc.called {
				t.Fatalf("service must not be called for invalid id %q", tc.id)
			}
		})
	}
}
