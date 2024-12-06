package collectors

import (
	"log/slog"
	"net/url"
	"strings"
	"sync"

	"github.com/DavidMovas/Movies-Reviews/internal/maps"

	"github.com/PuerkitoBio/goquery"

	"github.com/DavidMovas/Movies-Reviews/scraper/models"
	"github.com/gocolly/colly/v2"
)

type CastCollector struct {
	c *colly.Collector
	l *slog.Logger

	castMap   map[string]*models.Cast
	starLinks map[string]bool
	mx        sync.RWMutex
}

func NewCastCollector(c *colly.Collector, starCollector *StarCollector, logger *slog.Logger) *CastCollector {
	collector := &CastCollector{
		c:         c,
		l:         logger.With("collector", "cast"),
		castMap:   make(map[string]*models.Cast),
		starLinks: make(map[string]bool),
	}

	c.OnHTML("html", func(e *colly.HTMLElement) {
		movieID := getMovieID(e.Request.URL)

		cast := collector.getOrCreateCast(movieID, e.Request.URL.String())

		directorHeader := e.DOM.Find("h4#director")
		if directorHeader.Nodes != nil {
			collector.addCastFromSimpleTable(cast, "director", directorHeader.Next())
		}

		castHeader := e.DOM.Find("h4#cast")
		if castHeader.Nodes != nil {
			collector.addCastFromCastTable(cast, castHeader.Next(), 15)
		}

		writerHeader := e.DOM.Find("h4#writer")
		if writerHeader.Nodes != nil {
			collector.addCastFromSimpleTable(cast, "writer", writerHeader.Next())
		}

		producerHeader := e.DOM.Find("h4#producer")
		if producerHeader.Nodes != nil {
			collector.addCastFromSimpleTable(cast, "producer", producerHeader.Next())
		}

		composerHeader := e.DOM.Find("h4#composer")
		if composerHeader.Nodes != nil {
			collector.addCastFromSimpleTable(cast, "composer", composerHeader.Next())
		}

		collector.l.
			With("movie_id", movieID).
			Debug("cast collected")

		for _, credit := range cast.Cast {
			if collector.isNewStarLink(credit.StarLink) {
				starCollector.Visit(credit.StarLink)
			}
		}
	})

	return collector
}

func (c *CastCollector) Visit(link string) {
	visit(c.c, link, c.l)
}

func (c *CastCollector) Wait() {
	c.c.Wait()
}

func (c *CastCollector) Cast() map[string]*models.Cast {
	return c.castMap
}

func (c *CastCollector) getOrCreateCast(movieID, link string) *models.Cast {
	cast, _, _ := maps.GetOrCreateLocked(c.castMap, movieID, &c.mx, func(key string) (*models.Cast, error) {
		return &models.Cast{
			MovieID: key,
			Link:    link,
		}, nil
	})

	return cast
}

func (c *CastCollector) isNewStarLink(link string) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.starLinks[link]; ok {
		return false
	}

	c.starLinks[link] = true
	return true
}

func (c *CastCollector) addCastFromSimpleTable(cast *models.Cast, role string, table *goquery.Selection) {
	table.Find("tr").Each(func(_ int, row *goquery.Selection) {
		starLink := row.Find("td.name a")
		heroBlock := row.Find("a.cast-item-characters-link data-testid.cast-item-characters-link")
		heroName := heroBlock.Find("span")
		if starLink.Nodes == nil {
			return
		}

		href := starLink.AttrOr("href", "")
		link, _ := url.JoinPath("https://www.imdb.com", href)
		link = removeQueryPart(link)
		details := row.Find("td.credit").Text()

		credit := &models.Credit{
			Role:     role,
			Details:  strings.TrimSpace(details),
			StarName: strings.TrimSpace(starLink.Text()),
			HeroName: strings.TrimSpace(heroName.Text()),
			StarLink: link,
			StarID:   getStarID(link),
		}

		cast.Cast = append(cast.Cast, credit)
	})
}

func (c *CastCollector) addCastFromCastTable(cast *models.Cast, table *goquery.Selection, max int) {
	var added int
	table.Find("tr").EachWithBreak(func(_ int, row *goquery.Selection) bool {
		starLink := row.Find("td:not(.primary_photo .character) a")
		heroBlock := row.Find("a.cast-item-characters-link data-testid.cast-item-characters-link")
		heroName := heroBlock.Find("span")

		if starLink.Nodes == nil {
			return true
		}

		href := starLink.AttrOr("href", "")
		link, _ := url.JoinPath("https://www.imdb.com", href)
		link = removeQueryPart(link)

		credit := &models.Credit{
			Role:     "actor",
			Details:  "",
			StarName: strings.TrimSpace(starLink.Text()),
			StarLink: link,
			HeroName: strings.TrimSpace(heroName.Text()),
			StarID:   getStarID(link),
		}

		cast.Cast = append(cast.Cast, credit)
		added++

		return added <= max
	})
}
