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

type mockAuthService struct {
	resp *models.LoginResponse
	err  error

	called bool
}

func (m *mockAuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	m.called = true
	return m.resp, m.err
}

func TestPostLogin_Handler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// prepare mock auth service
	resp := &models.LoginResponse{ID: 1, Name: "Admin Lelang", Username: "admin", Role: models.RoleAdmin, Token: "a.b.c"}
	mockSvc := &mockAuthService{resp: resp}

	h := handlers.New(mockSvc, nil, nil, nil, nil, nil, nil)

	body := `{"username":"admin","password":"admin123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	h.PostLogin(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body=%s", http.StatusOK, w.Code, w.Body.String())
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
		t.Fatalf("expected login data object, body=%s", w.Body.String())
	}
	if data["token"] != "a.b.c" || data["username"] != "admin" {
		t.Fatalf("unexpected login payload: %+v", data)
	}
}

func TestPostLogin_Handler_ErrorMapping(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		svcErr   error
		wantCode int
	}{
		{"invalid credentials", svc.ErrInvalidCredentials, http.StatusUnauthorized},
		{"unexpected error", errors.New("db fail"), http.StatusInternalServerError},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAuthService{resp: nil, err: tc.svcErr}
			h := handlers.New(mockSvc, nil, nil, nil, nil, nil, nil)

			body := `{"username":"admin","password":"wrong"}`
			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h.PostLogin(c)

			if w.Code != tc.wantCode {
				t.Fatalf("expected status %d, got %d, body=%s", tc.wantCode, w.Code, w.Body.String())
			}
			if strings.Contains(w.Body.String(), "db fail") {
				t.Fatalf("internal error detail must not leak to the client, body=%s", w.Body.String())
			}
		})
	}
}

func TestPostLogin_Handler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		body        string
		wantMessage string
	}{
		{"missing password", `{"username":"admin"}`, "Missing or invalid field: password"},
		{"missing both fields", `{}`, "Missing or invalid field: username, password"},
		{"syntax error", `{"username":}`, "Malformed JSON body"},
		{"wrong field type", `{"username":123,"password":"x"}`, "Malformed JSON body"},
		{"truncated body", `{"username":`, "Invalid request body"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &mockAuthService{}
			h := handlers.New(mockSvc, nil, nil, nil, nil, nil, nil)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			h.PostLogin(c)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d, body=%s", http.StatusBadRequest, w.Code, w.Body.String())
			}

			var ar models.APIResponse
			if err := json.Unmarshal(w.Body.Bytes(), &ar); err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}
			if ar.Message != tc.wantMessage {
				t.Fatalf("expected message %q, got %q", tc.wantMessage, ar.Message)
			}
			if mockSvc.called {
				t.Fatalf("service must not be called when the body is invalid")
			}
		})
	}
}
