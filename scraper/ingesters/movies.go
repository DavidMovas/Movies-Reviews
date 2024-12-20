package ingesters

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/maps"
	"github.com/DavidMovas/Movies-Reviews/internal/slices"
	"golang.org/x/sync/errgroup"

	"github.com/DavidMovas/Movies-Reviews/contracts"

	"github.com/DavidMovas/Movies-Reviews/scraper/models"

	"github.com/DavidMovas/Movies-Reviews/client"
)

type MoviesIngester struct {
	c                *client.Client
	token            string
	genreIDConverter func(string) (int, bool)
	starIDConverter  func(string) (int, bool)
	logger           *slog.Logger
}

func NewMovieIngester(c *client.Client, token string, genreIngesterConverter func(string) (int, bool), starIngesterConverter func(string) (int, bool), logger *slog.Logger) *MoviesIngester {
	return &MoviesIngester{
		c:                c,
		token:            token,
		logger:           logger.With("ingester", "movies"),
		genreIDConverter: genreIngesterConverter,
		starIDConverter:  starIngesterConverter,
	}
}

func (i *MoviesIngester) Ingest(movies map[string]*models.Movie, casts map[string]*models.Cast) error {
	existingMovies, err := client.Paginate[*contracts.Movie](&contracts.GetMoviesRequest{}, i.c.GetMovies)
	if err != nil {
		return err
	}

	type movieCommonIdentifier struct {
		Title       string
		ReleaseDate time.Time
	}

	getIDFn := func(m *contracts.Movie) movieCommonIdentifier {
		return movieCommonIdentifier{
			Title:       m.Title,
			ReleaseDate: m.ReleaseDate,
		}
	}

	idToMovieMap := slices.ToMap(existingMovies, getIDFn, slices.NoChangeFunc[*contracts.Movie]())
	var mx sync.RWMutex

	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(8)

	for _, movie := range movies {
		movie := movie
		commonID := movieCommonIdentifier{
			Title:       movie.Title,
			ReleaseDate: movie.ReleaseDate,
		}

		if maps.ExistsLocked(idToMovieMap, commonID, &mx) {
			continue
		}

		group.Go(func() error {
			var created bool
			_, created, err = maps.GetOrCreateLocked(idToMovieMap, commonID, &mx, func(_ movieCommonIdentifier) (*contracts.Movie, error) {
				req := &contracts.CreateMovieRequest{
					Title:       movie.Title,
					ReleaseDate: movie.ReleaseDate,
					Description: movie.Description,
					IMDbURL:     &movie.Link,
				}

				if movie.PosterURL != "" {
					req.PosterURL = &movie.PosterURL
				}

				if movie.IMDbRating != 0 {
					req.IMDbRating = &movie.IMDbRating
				}

				if movie.Metascore != 0 {
					req.Metascore = &movie.Metascore
				}

				if movie.MetascoreURL != "" {
					req.MetascoreURL = &movie.MetascoreURL
				}

				for _, genre := range movie.Genres {
					genreID, ok := i.genreIDConverter(genre)
					if !ok {
						i.logger.
							With("genre", genre).
							Error("genre not found")
					}

					req.GenreIDs = append(req.GenreIDs, genreID)
				}

				cast, ok := casts[movie.ID]
				if !ok {
					i.logger.
						With("movieID", movie.ID).
						Error("cast not found")
				}

				for _, credit := range cast.Cast {
					starID, ok := i.starIDConverter(credit.StarID)
					if !ok {
						i.logger.
							With("starID", credit.StarID).
							Error("star not found")
						continue
					}

					creditInfo := &contracts.MovieCreditInfo{
						StarID:  starID,
						Role:    credit.Role,
						IMDbURL: &credit.StarLink,
					}

					if credit.Details != "" {
						creditInfo.Details = credit.Details
					}

					if credit.HeroName != "" {
						if len(credit.HeroName) < 100 {
							creditInfo.HeroName = &credit.HeroName
						}
					}

					req.Cast = append(req.Cast, *creditInfo)
				}

				var md *contracts.MovieDetails
				md, err = i.c.CreateMovie(contracts.NewAuthenticated(req, i.token))
				if err != nil {
					return nil, fmt.Errorf("create movie: %w", err)
				}

				return &md.Movie, nil
			})
			if err != nil {
				return err
			}

			if created {
				i.logger.
					With("movie_id", movie.ID).
					With("movie_common_id", commonID).
					Debug("movie created")
			}

			return nil
		})
	}

	if err = group.Wait(); err != nil {
		return fmt.Errorf("ingest movies: %w", err)
	}

	i.logger.Info("Ingested movies successfully")
	return nil
}
