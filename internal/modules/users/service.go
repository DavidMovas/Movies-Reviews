package users

import (
	"context"

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
	return s.repo.Create(ctx, user)
}

func (s *Service) GetExistingUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetExistingUserByEmail(ctx, email)
}

func (s *Service) GetExistingUserById(ctx context.Context, userId int) (*UserWithPassword, error) {
	return s.repo.GetExistingUserById(ctx, userId)
}

func (s *Service) GetExistingUserByUsername(ctx context.Context, username string) (*UserWithPassword, error) {
	return s.repo.GetExistingUserByUsername(ctx, username)
}

func (s *Service) UpdateExistingUserById(ctx context.Context, id int, user *NewUserData) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdateExistingUserById(ctx, id, user.Username, string(passHash))
}

func (s *Service) UpdateUserRoleById(ctx context.Context, id int, role string) error {
	return s.repo.UpdateUserRoleById(ctx, id, role)
}

func (s *Service) DeleteExistingUserById(ctx context.Context, userId int) error {
	return s.repo.DeleteExistingUserById(ctx, userId)
}

func (s *Service) CheckUserExistsById(ctx context.Context, id int) (bool, error) {
	return s.repo.CheckIsUserExistsById(ctx, id)
}
