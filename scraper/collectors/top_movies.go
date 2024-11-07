package collectors

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

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

	c.OnResponse(func(r *colly.Response) {
		time.Sleep(time.Second * 2)
		contentType := r.Headers.Get("Content-Type")
		collector.l.With("content_type", contentType).Debug("content type")
		if strings.Contains(contentType, "application/json") {
			var data map[string]any

			if err := json.Unmarshal(r.Body, &data); err != nil {
				collector.l.
					With("err", err).
					Error("failed to unmarshal top movies response")
				return
			}

			if id := findIDsWithPrefix(data, "tt"); id != "" {
				movieCollector.Visit("https://www.imdb.com/title/" + id)
			}
		}
	})

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

func findIDsWithPrefix(data map[string]interface{}, prefix string) string {
	for _, value := range data {
		switch v := value.(type) {
		case string:
			if strings.HasPrefix(v, prefix) {
				return v
			}
		case map[string]interface{}:
			findIDsWithPrefix(v, prefix)
		case []interface{}:
			for _, item := range v {
				if itemMap, ok := item.(map[string]interface{}); ok {
					findIDsWithPrefix(itemMap, prefix)
				}
			}
		}
	}

	return ""
}
