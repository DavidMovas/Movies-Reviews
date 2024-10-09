package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
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

func (s *Service) GenerateToken(claims *AccessClaims) (token string, err error) {
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GetAccessExpiration() time.Duration {
	return s.accessExpiration
}
