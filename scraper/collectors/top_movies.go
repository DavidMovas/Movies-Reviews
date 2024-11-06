package collectors

import (
	"log/slog"
	"strings"

	"github.com/gocolly/colly/v2"
)

type TopMoviesCollector struct {
	c *colly.Collector
	l *slog.Logger
}

func NewTopMoviesCollector(c *colly.Collector, movieCollector *MovieCollector, logger *slog.Logger) *TopMoviesCollector {
	_ = c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5})
	c.MaxDepth = 0

	collector := &TopMoviesCollector{
		c: c,
		l: logger.With("collector", "top_movies"),
	}

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, e *colly.HTMLElement) {
			link := e.Attr("href")
			text := strings.TrimSpace(e.Text)

			if text != "" && strings.HasPrefix(link, "/title/") {
				movieCollector.Visit(e.Request.AbsoluteURL(link))
			}
		})
	})

	return collector
}

func (c *TopMoviesCollector) Start() {
	starLink := "https://www.imdb.com/chart/top"
	visit(c.c, starLink, c.l)
}

func (c *TopMoviesCollector) Wait() {
	c.c.Wait()
}
