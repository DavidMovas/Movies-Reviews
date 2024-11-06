package collectors

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/gocolly/colly/v2"
)

func NewBaseCollector() *colly.Collector {
	baseCollector := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com"),
		colly.Async(true),
		colly.CacheDir("./.cache"),
	)

	return baseCollector
}

func Derive(c *colly.Collector) *colly.Collector {
	r := c.Clone()
	r.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
	})

	return r
}

func visit(c *colly.Collector, link string, logger *slog.Logger) {
	err := c.Visit(link)
	if err != nil {
		logger.
			With("link", link).
			With("error", err).
			Error("error visiting link")
	}
}

func removeQueryPart(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err).Error()
	}

	u.Path = strings.Split(u.Path, "/?")[0]
	u.RawQuery = ""
	return u.String()
}
