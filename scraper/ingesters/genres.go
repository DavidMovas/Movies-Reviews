package ingesters

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/DavidMovas/Movies-Reviews/internal/maps"

	"golang.org/x/sync/errgroup"

	"github.com/DavidMovas/Movies-Reviews/internal/slices"

	"github.com/DavidMovas/Movies-Reviews/contracts"

	"github.com/DavidMovas/Movies-Reviews/client"
)

type GenreIngester struct {
	c      *client.Client
	token  string
	logger *slog.Logger

	conversionMap map[string]int
}

func NewGenreIngest(c *client.Client, token string, logger *slog.Logger) *GenreIngester {
	return &GenreIngester{
		c:      c,
		token:  token,
		logger: logger,
	}
}

func (i *GenreIngester) Ingest(genres []string) error {
	existingGenres, err := i.c.GetGenres()
	if err != nil {
		return fmt.Errorf("get genres: %w", err)
	}

	getGenreNameFn := func(g *contracts.Genre) string {
		return g.Name
	}

	nameToGenreMap := slices.ToMap(existingGenres, getGenreNameFn, slices.NoChangeFunc[*contracts.Genre]())
	var mx sync.RWMutex

	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(8)

	for _, genre := range genres {
		genre := genre

		if maps.ExistsLocked(nameToGenreMap, genre, &mx) {
			continue
		}

		group.Go(func() error {
			var created bool
			_, created, err = maps.GetOrCreateLocked(nameToGenreMap, genre, &mx, func(name string) (*contracts.Genre, error) {
				req := &contracts.CreateGenreRequest{
					Name: name,
				}
				return i.c.CreateGenre(contracts.NewAuthenticated[contracts.CreateGenreRequest](*req, i.token))
			})
			if err != nil {
				return fmt.Errorf("create genre %q: %w", genre, err)
			}

			if created {
				i.logger.
					With("genre", genre).
					Debug("created genre")
			}

			return nil
		})
	}

	if err = group.Wait(); err != nil {
		return fmt.Errorf("ingest genres: %w", err)
	}

	i.conversionMap = make(map[string]int, len(nameToGenreMap))
	for _, genre := range nameToGenreMap {
		i.conversionMap[genre.Name] = genre.ID
	}

	i.logger.Info("Ingested genres successfully")
	return nil
}

func (i *GenreIngester) Converter(genre string) (int, bool) {
	id, ok := i.conversionMap[genre]
	return id, ok
}
