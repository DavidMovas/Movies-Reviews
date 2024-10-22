package movies

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"

	"github.com/DavidMovas/Movies-Reviews/internal/log"
)

type Service struct {
	repo       *Repository
	genresRepo *genres.Repository
}

func NewService(repo *Repository, genresRepo *genres.Repository) *Service {
	return &Service{
		repo:       repo,
		genresRepo: genresRepo,
	}
}

func (s *Service) GetMovies(ctx context.Context, offset int, limit int, sort, order string) ([]*Movie, int, error) {
	return s.repo.GetMovies(ctx, offset, limit, sort, order)
}

func (s *Service) GetMovieByID(ctx context.Context, movieID int) (*MovieDetails, error) {
	movie, err := s.repo.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}

	err = s.assemble(ctx, movie)

	return movie, err
}

func (s *Service) CreateMovie(ctx context.Context, movie *MovieDetails) (*MovieDetails, error) {
	err := s.repo.CreateMovie(ctx, movie)
	if err != nil {
		return nil, err
	}

	err = s.assemble(ctx, movie)

	log.FromContext(ctx).Info("movie created", "movie_id", movie.ID)
	return movie, err
}

func (s *Service) UpdateMovieByID(ctx context.Context, movieID int, req *UpdateMovieRequest) (*MovieDetails, error) {
	movie, err := s.repo.UpdateMovieByID(ctx, movieID, req)
	if err != nil {
		return nil, err
	}

	err = s.assemble(ctx, movie)

	log.FromContext(ctx).Info("movie updated", "movie_id", movie.ID)
	return movie, err
}

func (s *Service) DeleteMovieByID(ctx context.Context, movieID int) error {
	if err := s.repo.DeleteMovieByID(ctx, movieID); err != nil {
		return err
	}

	log.FromContext(ctx).Info("movie deleted", "movie_id", movieID)
	return nil
}

func (s *Service) assemble(ctx context.Context, movie *MovieDetails) error {
	var err error
	movie.Genres, err = s.genresRepo.GetGenresByMovieID(ctx, movie.ID)
	return err
}
