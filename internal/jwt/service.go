package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Service struct {
	secret           string
	accessExpiration time.Duration
}

func NewService(secret string, accessExpiration time.Duration) *Service {
	return &Service{
		secret:           secret,
		accessExpiration: accessExpiration,
	}
}

func (s *Service) GenerateToken(userID int, role string) (string, error) {
	now := time.Now()
	claims := &AccessClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			Subject:   strconv.Itoa(userID),
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(s.accessExpiration).Unix(),
		},
		UserID: userID,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

func (s *Service) GetAccessExpiration() time.Duration {
	return s.accessExpiration
}
