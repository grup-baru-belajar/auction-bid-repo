package services

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/grup-baru-belajar/auction-bid-repo/internal/models"
)

type Claims struct {
	Username string `json:"username"`
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
	ttl time.Duration
}

func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

func (m *TokenManager) Generate(user *models.User) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: user.Username,
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.FormatInt(user.ID, 10),
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}