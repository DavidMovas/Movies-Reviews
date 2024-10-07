package jwt

import "time"

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

func (s *Service) GenerateToken(id int, role string) (token string, err error) {

	//TODO:
	// 1. Create token from claims using secret
	// 2. Return token

	return "", nil
}
