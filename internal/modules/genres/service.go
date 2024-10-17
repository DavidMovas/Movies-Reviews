package genres

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/log"
)

type Service struct {
	Repository *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		Repository: repo,
	}
}

func (s *Service) CreateGenre(ctx context.Context, raq *contracts.CreateGenreRequest) (*contracts.Genre, error) {
	genre, err := s.Repository.CreateGenre(ctx, raq)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("genre created", "genre_id", genre.ID)
	return genre, nil
}

func (s *Service) GetGenres(context context.Context) ([]*contracts.Genre, error) {
	return s.Repository.GetGenres(context)
}

func (s *Service) GetGenreById(ctx context.Context, id int) (*contracts.Genre, error) {
	return s.Repository.GetGenreById(ctx, id)
}

func (s *Service) UpdateGenreById(ctx context.Context, id int, raq *contracts.UpdateGenreRequest) error {
	if err := s.Repository.UpdateGenreById(ctx, id, raq); err != nil {
		return err
	}

	log.FromContext(ctx).Info("genre updated", "genre_id", id)
	return nil
}

func (s *Service) DeleteGenreById(ctx context.Context, id int) error {
	if err := s.Repository.DeleteGenreById(ctx, id); err != nil {
		return err
	}

	log.FromContext(ctx).Info("genre deleted", "genre_id", id)
	return nil
}
