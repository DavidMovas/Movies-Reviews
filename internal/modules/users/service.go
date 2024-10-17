package users

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/contracts"
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

func (s *Service) Create(ctx context.Context, user *contracts.UserWithPassword) error {
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user created", "user_id", user.ID)
	return nil
}

func (s *Service) GetExistingUserByEmail(ctx context.Context, email string) (*contracts.UserWithPassword, error) {
	return s.repo.GetExistingUserByEmail(ctx, email)
}

func (s *Service) GetExistingUserById(ctx context.Context, userId int) (*contracts.UserWithPassword, error) {
	return s.repo.GetExistingUserById(ctx, userId)
}

func (s *Service) GetExistingUserByUsername(ctx context.Context, username string) (*contracts.UserWithPassword, error) {
	return s.repo.GetExistingUserByUsername(ctx, username)
}

func (s *Service) UpdateExistingUserById(ctx context.Context, id int, user *contracts.UpdateUserRequest) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repo.UpdateExistingUserById(ctx, id, user.Username, string(passHash)); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user updated", "user_id", id)
	return nil
}

func (s *Service) UpdateUserRoleById(ctx context.Context, id int, role string) error {
	if err := s.repo.UpdateUserRoleById(ctx, id, role); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user role updated", "user_id", id, "role", role)
	return nil
}

func (s *Service) DeleteExistingUserById(ctx context.Context, userId int) error {
	if err := s.repo.DeleteExistingUserById(ctx, userId); err != nil {
		return err
	}

	log.FromContext(ctx).Info("user deleted", "user_id", userId)
	return nil
}
