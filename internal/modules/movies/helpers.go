package movies

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/stars"
	"github.com/DavidMovas/Movies-Reviews/internal/slices"
	"golang.org/x/sync/errgroup"
)

func (s *Service) assemble(ctx context.Context, movie *MovieDetails) error {
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		var err error
		movie.Genres, err = s.genresRepo.GetGenresByMovieID(groupCtx, movie.ID)
		return err
	})

	group.Go(func() error {
		var err error
		var credits []*stars.MovieCredit
		credits, err = s.starsRepo.GetStarsByMovieID(groupCtx, movie.ID)
		if err == nil {
			movie.Cast = slices.CastSlice(credits, func(credit *stars.MovieCredit) *MovieCredit {
				return &MovieCredit{
					Star:     credit.Star,
					HeroName: credit.HeroName,
					Role:     credit.Role,
					Details:  credit.Details,
				}
			})
		}

		return err
	})

	return group.Wait()
}
