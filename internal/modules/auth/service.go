package auth

import (
	"context"
	"errors"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	usersService *users.Service
}

func (s *Service) Register(ctx context.Context, user *users.User, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userWithPassword := &users.UserWithPassword{
		User:         *user,
		PasswordHash: string(passHash),
	}

	return s.usersService.Create(ctx, userWithPassword)
}

func (s *Service) Login(ctx context.Context, email, password string) (token string, err error) {
	// TODO:
	// 1. Get user by email from DB. Using userService
	// 2. bcrypt.CompareHashAndPassword() need to compare hash and password
	// 3. Create Claims from user data
	// 4. Create JWT token

	return "", errors.New("not implemented")
}

func NewService(service *users.Service) *Service {
	return &Service{
		usersService: service,
	}
}
