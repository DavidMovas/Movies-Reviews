package collectors

import (
	"encoding/json"
	"github.com/gocolly/colly/v2"
	"log/slog"
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

	c.OnHTML("script[type='application/ld+json']", func(e *colly.HTMLElement) {
		var data map[string]any

		jsonText := e.Text

		if err := json.Unmarshal([]byte(jsonText), &data); err != nil {
			collector.l.With("err", err).Error("error unmarshalling data")
			return
		}

		urls := findMovieURLs(data)
		for _, url := range urls {
			movieCollector.Visit(url)
		}
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

func findMovieURLs(data map[string]interface{}) []string {
	var urls []string

	for key, value := range data {
		switch v := value.(type) {
		case string:
			if key == "url" && isMovieURL(v) {
				urls = append(urls, v)
			}
		case map[string]interface{}:
			urls = append(urls, findMovieURLs(v)...)
		case []interface{}:
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					urls = append(urls, findMovieURLs(itemMap)...)
				}
			}
		}
	}

	return urls
}

func isMovieURL(url string) bool {
	return len(url) > 0 && url[:29] == "https://www.imdb.com/title/tt"
}
