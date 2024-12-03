package auth

import (
	"context"
	"errors"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	usersService *users.Service
	jwtService   *jwt.Service
}

func NewService(service *users.Service, jwtService *jwt.Service) *Service {
	return &Service{
		usersService: service,
		jwtService:   jwtService,
	}
}

func (s *Service) Register(ctx context.Context, user *users.User, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Internal(err)
	}

	userWithPassword := &users.UserWithPassword{
		User:         user,
		PasswordHash: string(passHash),
	}

	return s.usersService.Create(ctx, userWithPassword)
}

func (s *Service) Login(ctx context.Context, email, username *string, password string) (user *users.UserWithPassword, token string, err error) {
	if email != nil {
		user, err = s.usersService.GetExistingUserByEmail(ctx, *email)
	} else {
		user, err = s.usersService.GetExistingUserByUsername(ctx, *username)
	}

	if err != nil {
		return nil, "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, "", apperrors.Unauthorized("invalid password")
		}
		return nil, "", apperrors.Internal(err)
	}

	token, err = s.jwtService.GenerateToken(user.ID, user.Role)
	return user, token, err
}
