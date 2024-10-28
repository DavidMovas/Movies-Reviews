package stars

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/internal/log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetStars(ctx context.Context) ([]*Star, error) {
	return s.repo.GetStars(ctx)
}

func (s *Service) GetStarByID(ctx context.Context, starID int) (*Star, error) {
	return s.repo.GetStarByID(ctx, starID)
}

func (s *Service) GetStarsPaginated(ctx context.Context, offset int, limit int) ([]*Star, int, error) {
	return s.repo.GetStarsPaginated(ctx, offset, limit)
}

func (s *Service) CreateStar(ctx context.Context, req *CreateStarRequest) (*Star, error) {
	star, err := s.repo.CreateStar(ctx, req)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("star created", "star_id", star.ID)
	return star, nil
}

func (s *Service) UpdateStar(ctx context.Context, starID int, req *UpdateStarRequest) (*Star, error) {
	star, err := s.repo.UpdateStar(ctx, starID, req)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("star updated", "star_id", star.ID)
	return star, nil
}

func (s *Service) DeleteStarByID(ctx context.Context, starID int) error {
	if err := s.repo.DeleteStarByID(ctx, starID); err != nil {
		return err
	}

	log.FromContext(ctx).Info("star deleted", "star_id", starID)
	return nil
}
