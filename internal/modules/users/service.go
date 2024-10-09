package users

import "context"

type Service struct {
	repo *Repository
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	return s.repo.Create(ctx, user)
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}
