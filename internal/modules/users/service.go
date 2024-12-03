package users

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/internal/log"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user created", "user_id", user.ID)
	return nil
}

func (s *Service) GetExistingUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetExistingUserByEmail(ctx, email)
}

func (s *Service) GetExistingUserByID(ctx context.Context, userID int) (*User, error) {
	return s.repo.GetExistingUserByID(ctx, userID)
}

func (s *Service) GetExistingUserByUsername(ctx context.Context, username string) (*UserWithPassword, error) {
	return s.repo.GetExistingUserByUsername(ctx, username)
}

func (s *Service) UpdateExistingUserByID(ctx context.Context, userID int, user *UpdateUserRequest) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err = s.repo.UpdateExistingUserByID(ctx, userID, user.Username, string(passHash)); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user updated", "user_id", userID)
	return nil
}

func (s *Service) UpdateUserRoleByID(ctx context.Context, userID int, role string) error {
	if err := s.repo.UpdateUserRoleByID(ctx, userID, role); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user role updated", "user_id", userID, "role", role)
	return nil
}

func (s *Service) DeleteExistingUserByID(ctx context.Context, userID int) error {
	if err := s.repo.DeleteExistingUserByID(ctx, userID); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user deleted", "user_id", userID)
	return nil
}
