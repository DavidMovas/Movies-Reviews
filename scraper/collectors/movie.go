package collectors

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/maps"

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

		ratingTable := e.DOM.Find("div.sc-d541859f-2")
		if ratingTable.Nodes != nil {
			collector.getIMDbRating(movie, ratingTable)
		}

		metascoreHeader := e.DOM.Find("a.isReview")
		if metascoreHeader.Nodes != nil {
			collector.getMetascore(movie, metascoreHeader)
		}

		movie.Title = info.Name
		movie.Description = info.Description
		movie.Genres = info.Genre
		movie.ReleaseDate = mustParseDate(info.DatePublished)
		movie.PosterURL = info.Image
		movie.MetascoreURL, _ = url.JoinPath(info.URL, "/criticreviews")

		collector.toAllGenres(movie.Genres)

		collector.l.
			With("movie_id", movieID).
			With("title", movie.Title).
			Debug("movie collected")

		//creditsLink, _ := url.JoinPath(info.URL, "/fullcredits")
		//castCollector.Visit(creditsLink)
	})

	return collector
}

func (c *MovieCollector) Visit(link string) {
	visit(c.c, link, c.l)
}

func (c *MovieCollector) Wait() {
	c.c.Wait()
}

func (c *MovieCollector) Movies() map[string]*models.Movie {
	return c.movieMap
}

func (c *MovieCollector) Genres() []string {
	genres := make([]string, 0, len(c.allGenres))
	for genre := range c.allGenres {
		genres = append(genres, genre)
	}

	return genres
}

func (c *MovieCollector) getOrCreateMovie(movieID, link string) *models.Movie {
	movie, _, _ := maps.GetOrCreateLocked(c.movieMap, movieID, &c.mx, func(_ string) (*models.Movie, error) {
		return &models.Movie{
			ID:   movieID,
			Link: removeQueryPart(link),
		}, nil
	})

	return movie
}

func (c *MovieCollector) toAllGenres(genres []string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	for _, genre := range genres {
		c.allGenres[genre] = true
	}
}

func (c *MovieCollector) getIMDbRating(movie *models.Movie, table *goquery.Selection) {
	table.Find("span.sc-d541859f-1").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		scoreLine := row.Text()

		if scoreLine != "" {
			var num float64
			num, _ = strconv.ParseFloat(scoreLine, 32)
			movie.IMDbRating = math.Round(num*10) / 10
			return false
		}

		return true
	})
}

func (c *MovieCollector) getMetascore(movie *models.Movie, table *goquery.Selection) {
	scoreBox := table.Find("span.metacritic-score-box")
	movie.Metascore, _ = strconv.Atoi(scoreBox.Text())
}

func getMovieID(url *url.URL) string {
	id := strings.Split(url.Path, "/")[2]
	return id
}

func mustParseDate(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}
	}

	return t
}
