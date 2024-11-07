package ingesters

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gopkg.in/validator.v2"

	"github.com/DavidMovas/Movies-Reviews/internal/maps"

	"github.com/DavidMovas/Movies-Reviews/internal/slices"
	"golang.org/x/sync/errgroup"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/scraper/models"

	"github.com/DavidMovas/Movies-Reviews/client"
)

type StarIngester struct {
	c      *client.Client
	token  string
	logger *slog.Logger

	conversionMap map[string]int
}

func NewStarIngester(c *client.Client, token string, logger *slog.Logger) *StarIngester {
	return &StarIngester{
		c:      c,
		token:  token,
		logger: logger.With("ingester", "stars"),
	}
}

func (i *StarIngester) Ingest(stars map[string]*models.Star, bios map[string]*models.Bio) error {
	existingStars, err := client.Paginate(&contracts.GetStarsRequest{}, i.c.GetStars)
	if err != nil {
		return fmt.Errorf("get stars: %w", err)
	}

	type starCommonIdentifier struct {
		FirstName string `validate:"min=3,max=50"`
		Lastname  string `validate:"min=3,max=50"`
		BirthDate time.Time
	}

	getIDFn := func(s *contracts.Star) starCommonIdentifier {
		return starCommonIdentifier{
			FirstName: s.FirstName,
			Lastname:  s.LastName,
			BirthDate: s.BirthDate,
		}
	}

	idToStarMap := slices.ToMap(existingStars, getIDFn, slices.NoChangeFunc[*contracts.Star]())
	var mx sync.RWMutex

	group, _ := errgroup.WithContext(context.Background())
	group.SetLimit(8)

	for _, star := range stars {
		star := star
		commonID := starCommonIdentifier{
			FirstName: star.FirstName,
			Lastname:  star.LastName,
			BirthDate: star.BirthDate,
		}

		if validator.Validate(commonID) != nil {
			continue
		}

		if maps.ExistsLocked(idToStarMap, commonID, &mx) {
			continue
		}

		group.Go(func() error {
			var created bool
			_, created, err = maps.GetOrCreateLocked(idToStarMap, commonID, &mx, func(_ starCommonIdentifier) (*contracts.Star, error) {
				bio, ok := bios[star.ID]
				if !ok {
					i.logger.With("star_id", star.ID).
						Error("no bio found for star")
					bio = &models.Bio{}
				}

				req := &contracts.CreateStarRequest{
					FirstName: star.FirstName,
					LastName:  star.LastName,
					BirthDate: star.BirthDate,
					DeathDate: star.DeathDate,
				}

				if bio.Bio != "" {
					req.Bio = &bio.Bio
				}

				if bio.BirthPlace != "" {
					req.BirthPlace = &bio.BirthPlace
				}

				var s *contracts.Star
				s, err = i.c.CreateStar(contracts.NewAuthenticated(req, i.token))
				if err != nil {
					return nil, fmt.Errorf("create star: %w", err)
				}

				return s, nil
			})
			if err != nil {
				return err
			}

			if created {
				i.logger.
					With("star_id", star.ID).
					With("star_common_id", commonID).
					Debug("star created")
			}

			return nil
		})
	}

	if err = group.Wait(); err != nil {
		return fmt.Errorf("ingest stars: %w", err)
	}

	i.conversionMap = make(map[string]int, len(idToStarMap))
	for _, star := range stars {
		commonID := starCommonIdentifier{
			FirstName: star.FirstName,
			Lastname:  star.LastName,
			BirthDate: star.BirthDate,
		}
		s, ok := idToStarMap[commonID]
		if !ok {
			i.logger.With("star_id", star.ID).
				Error("no star found for id")
			continue
		}

		i.conversionMap[star.ID] = s.ID
	}

	i.logger.Info("Stars ingested successfully")
	return nil
}

func (i *StarIngester) Converter(imdbID string) (int, bool) {
	id, ok := i.conversionMap[imdbID]
	return id, ok
}
