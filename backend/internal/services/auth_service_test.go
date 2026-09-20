package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
	svc "github.com/grup-baru-belajar/auction-bid-repo/internal/services"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	result *models.User
	err    error

	gotUsername string
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	m.gotUsername = username
	if m.err != nil {
		return nil, m.err
	}
	return m.result, nil
}

func TestLogin_Success(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := &models.User{
		ID:       2,
		Name:     "John Doe",
		Username: "john",
		Password: string(hash),
		Role:     models.RoleUser,
	}

	mock := &mockUserRepo{result: user}
	tokens := token.NewTokenManager("test-secret-long-enough-for-hs256", time.Hour)
	s := svc.NewAuthService(mock, tokens)

	req := models.LoginRequest{Username: "john", Password: "password123"}
	got, err := s.Login(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != user.ID || got.Name != user.Name || got.Username != user.Username || got.Role != user.Role {
		t.Fatalf("unexpected login response: %+v", got)
	}
	if got.Token == "" {
		t.Fatalf("expected token to be set, got empty string")
	}
	if mock.gotUsername != "john" {
		t.Fatalf("expected repo called with john, got %q", mock.gotUsername)
	}
}

func TestLogin_TokenClaims(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := &models.User{ID: 2, Name: "John Doe", Username: "john", Password: string(hash), Role: models.RoleUser}

	tokens := token.NewTokenManager("test-secret-long-enough-for-hs256", time.Hour)
	s := svc.NewAuthService(&mockUserRepo{result: user}, tokens)

	got, err := s.Login(context.Background(), models.LoginRequest{Username: "john", Password: "password123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	claims, err := tokens.Parse(got.Token)
	if err != nil {
		t.Fatalf("issued token should be parseable: %v", err)
	}
	id, err := claims.UserID()
	if err != nil {
		t.Fatalf("subject should hold the user id: %v", err)
	}
	if id != user.ID || claims.Username != user.Username || claims.Role != user.Role {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestLogin_Errors(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := &models.User{ID: 2, Name: "John Doe", Username: "john", Password: string(hash), Role: models.RoleUser}

	tests := []struct {
		name     string
		repoErr  error
		password string
		wantErr  error
	}{
		{"user not found", repository.ErrUserNotFound, "password123", svc.ErrInvalidCredentials},
		{"wrong password", nil, "wrong-password", svc.ErrInvalidCredentials},
		{"other repo error", errors.New("db fail"), "password123", nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mock := &mockUserRepo{result: user, err: tc.repoErr}
			tokens := token.NewTokenManager("test-secret-long-enough-for-hs256", time.Hour)
			s := svc.NewAuthService(mock, tokens)

			req := models.LoginRequest{Username: "john", Password: tc.password}
			got, err := s.Login(context.Background(), req)
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
				if errors.Is(err, svc.ErrInvalidCredentials) {
					t.Fatalf("db failure must not be reported as invalid credentials")
				}
			}
		})
	}
}

func TestLogin_SameErrorForUnknownUserAndWrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := &models.User{ID: 2, Name: "John Doe", Username: "john", Password: string(hash), Role: models.RoleUser}
	tokens := token.NewTokenManager("test-secret-long-enough-for-hs256", time.Hour)

	unknown := svc.NewAuthService(&mockUserRepo{err: repository.ErrUserNotFound}, tokens)
	_, errUnknown := unknown.Login(context.Background(), models.LoginRequest{Username: "ghost", Password: "password123"})

	wrong := svc.NewAuthService(&mockUserRepo{result: user}, tokens)
	_, errWrong := wrong.Login(context.Background(), models.LoginRequest{Username: "john", Password: "wrong-password"})

	if errUnknown == nil || errWrong == nil {
		t.Fatalf("expected both cases to fail, got %v and %v", errUnknown, errWrong)
	}
	if errUnknown.Error() != errWrong.Error() {
		t.Fatalf("errors must be identical so the response cannot reveal which field was wrong: %q vs %q", errUnknown, errWrong)
	}
}
