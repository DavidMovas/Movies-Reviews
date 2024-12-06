package moviesv2

import (
	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/movies"
	starsv2 "github.com/DavidMovas/Movies-Reviews/internal/modules/stars/v2"
)

type MovieDetailsV2 struct {
	movies.Movie
	Description string           `json:"description"`
	Version     int              `json:"version"`
	Genres      []*genres.Genre  `json:"genres"`
	Cast        []*MovieCreditV2 `json:"cast"`
}

type MovieCreditV2 struct {
	Star    starsv2.StarShort `json:"star"`
	Role    string            `json:"role"`
	Details string            `json:"details,omitempty"`
}
