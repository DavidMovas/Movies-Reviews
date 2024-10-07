package users

import "context"

type Service struct {
	repo *Repository
}

func (s *Service) Create(ctx context.Context, user *UserWithPassword) error {
	return s.repo.Create(ctx, user)
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}
