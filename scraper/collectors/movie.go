package collectors

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/DavidMovas/Movies-Reviews/scraper/models"

	"github.com/gocolly/colly/v2"
)

type MovieCollector struct {
	c *colly.Collector
	l *slog.Logger

	movieMap  map[string]*models.Movie
	allGenres map[string]bool
	mx        sync.RWMutex
}

func NewMovieCollector(c *colly.Collector, castCollector *CastCollector, logger *slog.Logger) *MovieCollector {
	collector := &MovieCollector{
		c:         c,
		l:         logger.With("collector", "movie"),
		movieMap:  make(map[string]*models.Movie),
		allGenres: make(map[string]bool),
	}

	c.OnHTML("html", func(e *colly.HTMLElement) {
		movieID := getMovieID(e.Request.URL)

		movie := collector.getOrCreateMovie(movieID, e.Request.URL.String())

		type movieInfo struct {
			URL           string   `json:"url"`
			Name          string   `json:"name"`
			Image         string   `json:"image"`
			Description   string   `json:"description"`
			Genre         []string `json:"genre"`
			DatePublished string   `json:"datePublished"`
		}

		var info movieInfo
		err := json.Unmarshal([]byte(e.ChildText("script[type='application/ld+json']")), &info)
		if err != nil {
			collector.l.
				With("movie_id", movieID).
				With("err", err).
				Error("error unmarshalling movie info")

			return
		}

		movie.Title = info.Name
		movie.Description = info.Description
		movie.Genres = info.Genre
		movie.ReleaseDate = mustParseDate(info.DatePublished)

		collector.toAllGenres(movie.Genres)

		collector.l.
			With("movie_id", movieID).
			With("title", movie.Title)
	})

	return collector
}
