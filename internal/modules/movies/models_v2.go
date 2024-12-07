package movies

import (
	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	starsv2 "github.com/DavidMovas/Movies-Reviews/internal/modules/stars"
)

type MovieDetailsV2 struct {
	Movie
	Description  string           `json:"description"`
	IMDbRating   *float64         `json:"imdbRating,omitempty"`
	IMDbURL      *string          `json:"imdbUrl,omitempty"`
	Metascore    *int             `json:"metascore,omitempty"`
	MetascoreURL *string          `json:"metascoreUrl,omitempty"`
	Version      int              `json:"version"`
	Genres       []*genres.Genre  `json:"genres"`
	Cast         []*MovieCreditV2 `json:"cast"`
}

type MovieCreditV2 struct {
	Star     *starsv2.StarV2 `json:"star"`
	HeroName *string         `json:"heroName,omitempty"`
	Role     string          `json:"role"`
	Details  string          `json:"details,omitempty"`
}

func (m *MovieDetails) ConvertToV2() *MovieDetailsV2 {
	var cast []*MovieCreditV2
	for _, c := range m.Cast {
		cast = append(cast, c.ConvertToV2())
	}

	return &MovieDetailsV2{
		Movie:        m.Movie,
		IMDbRating:   m.IMDbRating,
		IMDbURL:      m.IMDbURL,
		Metascore:    m.Metascore,
		MetascoreURL: m.MetascoreURL,
		Description:  m.Description,
		Version:      m.Version,
		Genres:       m.Genres,
		Cast:         cast,
	}
}

func (c *MovieCredit) ConvertToV2() *MovieCreditV2 {
	return &MovieCreditV2{
		Star:     c.Star.ConvertToV2(),
		HeroName: c.HeroName,
		Role:     c.Role,
		Details:  c.Details,
	}
}
