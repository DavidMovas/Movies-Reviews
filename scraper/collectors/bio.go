package collectors

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/DavidMovas/Movies-Reviews/internal/maps"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"

	"github.com/DavidMovas/Movies-Reviews/scraper/models"
	"github.com/gocolly/colly/v2"
)

type BioCollector struct {
	c *colly.Collector
	l *slog.Logger

	bioMap map[string]*models.Bio
	mx     sync.RWMutex
}

func NewBioCollector(c *colly.Collector, logger *slog.Logger) *BioCollector {
	_ = c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 3})

	collector := &BioCollector{
		c:      c,
		l:      logger.With("collector", "bio"),
		bioMap: make(map[string]*models.Bio),
	}

	c.OnHTML("html", func(e *colly.HTMLElement) {
		starID := getStarID(e.Request.URL.String())

		bio := collector.getOrCreateBio(starID, e.Request.URL.String())

		overviewSection := getSectionByTitle(e, "overview")
		if overviewSection.Nodes != nil {
			overviewSection.Find("a").Each(func(_ int, selection *goquery.Selection) {
				if bio.BirthPlace != "" {
					return
				}

				href := selection.AttrOr("href", "")
				if strings.Contains(href, "?birth_place") {
					bio.BirthPlace = strings.TrimSpace(selection.Text())
				}
			})
		}

		bioSelection := getSectionByTitle(e, "mini_bio")
		if bioSelection.Nodes != nil {
			bioContentContainer := bioSelection.Find("ul li > div ul > div:not(.ipc-metadata-list-item-html-item--subtext) div")
			if bioContentContainer.Nodes != nil {
				bio.Bio = bioContentContainer.Text()
				bio.Bio = collectText(bioContentContainer)
			}
		}

		collector.l.
			With("star_id", starID).
			Debug("bio collected")
	})
	return collector
}

func (c *BioCollector) Visit(link string) {
	visit(c.c, link, c.l)
}

func (c *BioCollector) Wait() {
	c.c.Wait()
}

func (c *BioCollector) Bios() map[string]*models.Bio {
	return c.bioMap
}

func (c *BioCollector) getOrCreateBio(starID, link string) *models.Bio {
	bio, _, _ := maps.GetOrCreateLocked(c.bioMap, starID, &c.mx, func(_ string) (*models.Bio, error) {
		return &models.Bio{
			ID:   starID,
			Link: link,
		}, nil
	})

	return bio
}

func getSectionByTitle(e *colly.HTMLElement, title string) *goquery.Selection {
	selector := fmt.Sprintf("div[data-testid=sub-section-%s]", title)
	return e.DOM.Find(selector)
}

func collectText(s *goquery.Selection) string {
	var buf bytes.Buffer

	var f func(node *html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.TextNode:
			text := strings.ReplaceAll(n.Data, "\n", " ")
			buf.WriteString(text)

		case html.ElementNode:
			switch n.Data {
			case "br":
				prev := n.PrevSibling
				firstBr := prev == nil || prev.Type != html.ElementNode || prev.Data != "br"
				if firstBr {
					buf.WriteString("\n")
				}
			default:
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
		default:
		}
	}

	f(s.Nodes[0])

	return buf.String()
}
