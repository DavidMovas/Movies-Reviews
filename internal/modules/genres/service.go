package genres

import (
	"context"

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

func (s *Service) CreateGenre(ctx context.Context, raq *CreateGenreRequest) (*Genre, error) {
	genre, err := s.Repository.CreateGenre(ctx, raq)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("genre created", "genre_id", genre.ID)
	return genre, nil
}

func (s *Service) GetGenres(context context.Context) ([]*Genre, error) {
	return s.Repository.GetGenres(context)
}

func (s *Service) GetGenreByID(ctx context.Context, id int) (*Genre, error) {
	return s.Repository.GetGenreByID(ctx, id)
}

func (s *Service) UpdateGenreByID(ctx context.Context, id int, raq *UpdateGenreRequest) error {
	if err := s.Repository.UpdateGenreByID(ctx, id, raq); err != nil {
		return err
	}

	log.FromContext(ctx).Info("genre updated", "genre_id", id)
	return nil
}

func (s *Service) DeleteGenreByID(ctx context.Context, id int) error {
	if err := s.Repository.DeleteGenreByID(ctx, id); err != nil {
		return err
	}

	log.FromContext(ctx).Info("genre deleted", "genre_id", id)
	return nil
}
