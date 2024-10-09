package auth

import (
	"context"
	"fmt"

	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	usersService *users.Service
	jwtService   *jwt.Service
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
	userPass, err := s.usersService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userPass.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	user := users.User{
		ID:       userPass.ID,
		Role:     userPass.Role,
		Username: userPass.Username,
		Email:    userPass.Email,
	}

	userClaims := jwt.NewAccessClaimsFromUser(&user, s.jwtService.GetAccessExpiration())
	token, err = s.jwtService.GenerateToken(userClaims)

	if err != nil {
		return "", err
	}

	return token, nil
}

func NewService(service *users.Service, jwtService *jwt.Service) *Service {
	return &Service{
		usersService: service,
		jwtService:   jwtService,
	}
}
