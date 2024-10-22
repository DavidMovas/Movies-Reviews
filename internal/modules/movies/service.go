package movies

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

func (s *Service) GetMovies(ctx context.Context, offset int, limit int, sort, order string) ([]*Movie, int, error) {
	return s.repo.GetMovies(ctx, offset, limit, sort, order)
}

func (s *Service) GetMovieByID(ctx context.Context, movieID int) (*MovieDetails, error) {
	return s.repo.GetMovieByID(ctx, movieID)
}

func (s *Service) CreateMovie(ctx context.Context, movie *MovieDetails) (*MovieDetails, error) {
	movie, err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("movie created", "movie_id", movie.ID)
	return movie, nil
}

func (s *Service) UpdateMovieByID(ctx context.Context, movieID int, req *UpdateMovieRequest) (*MovieDetails, error) {
	movie, err := s.repo.UpdateMovieByID(ctx, movieID, req)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("movie updated", "movie_id", movie.ID)

	return movie, nil
}

func (s *Service) DeleteMovieByID(ctx context.Context, movieID int) error {
	if err := s.repo.DeleteMovieByID(ctx, movieID); err != nil {
		return err
	}

	log.FromContext(ctx).Info("movie deleted", "movie_id", movieID)
	return nil
}
