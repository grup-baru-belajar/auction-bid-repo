package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthService interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error)
}

type authService struct {
	users repository.UserRepository
	tokens *TokenManager
}

func NewAuthService(users repository.UserRepository, tokens *TokenManager) AuthService {
	return &authService{users: users, tokens: tokens}
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	user, err := s.users.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err , repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokens.Generate(user)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &models.LoginResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Role:     user.Role,
		Token:    token,
	}, nil

}