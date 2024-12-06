package movies

import (
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/stars"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"
)

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	PosterURL   *string    `json:"posterUrl,omitempty"`
	ReleaseDate time.Time  `json:"releaseDate"`
	AvgRating   *float64   `json:"avgRating,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type MovieDetails struct {
	Movie
	Description  string          `json:"description"`
	IMDbRating   *float64        `json:"imdbRating,omitempty"`
	IMDbURL      *string         `json:"imdbUrl,omitempty"`
	Metascore    *int            `json:"metascore,omitempty"`
	MetascoreURL *string         `json:"metascoreUrl,omitempty"`
	Version      int             `json:"version"`
	Genres       []*genres.Genre `json:"genres"`
	Cast         []*MovieCredit  `json:"cast"`
}

type MovieCredit struct {
	Star     stars.Star `json:"star"`
	Role     string     `json:"role"`
	HeroName *string    `json:"heroName,omitempty"`
	Details  string     `json:"details,omitempty"`
}

type MovieCreditInfo struct {
	StarID   int     `json:"starId"`
	Role     string  `json:"role"`
	HeroName *string `json:"heroName,omitempty"`
	Details  string  `json:"details,omitempty"`
}

type GetMovieRequest struct {
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}

type GetMoviesRequest struct {
	pagination.PaginatedRequestOrdered
	SearchTerm *string `query:"q"`
}

type CreateMovieRequest struct {
	Title        string            `json:"title" validate:"min=1,max=100"`
	ReleaseDate  time.Time         `json:"releaseDate" validate:"nonzero"`
	PosterURL    *string           `json:"posterUrl,omitempty"`
	IMDbRating   *float64          `json:"imdbRating,omitempty"`
	IMDbURL      *string           `json:"imdbUrl,omitempty"`
	Metascore    *int              `json:"metascore,omitempty"`
	MetascoreURL *string           `json:"metascoreUrl,omitempty"`
	Description  string            `json:"description"`
	GenreIDs     []int             `json:"genreIds" validate:"nonzero"`
	Cast         []MovieCreditInfo `json:"cast"`
}

type UpdateMovieRequest struct {
	MovieID      int                `json:"-" param:"movieId" validate:"nonzero"`
	Title        *string            `json:"title,omitempty" validate:"max=100"`
	ReleaseDate  *time.Time         `json:"releaseDate,omitempty"`
	PosterURL    *string            `json:"posterUrl,omitempty"`
	IMDbRating   *float64           `json:"imdbRating,omitempty"`
	IMDbURL      *string            `json:"imdbUrl,omitempty"`
	Metascore    *int               `json:"metascore,omitempty"`
	MetascoreURL *string            `json:"metascoreUrl,omitempty"`
	Description  *string            `json:"description,omitempty"`
	Version      int                `json:"version"`
	GenreIDs     []*int             `json:"genreIds,omitempty"`
	Cast         []*MovieCreditInfo `json:"cast,omitempty"`
}

type DeleteMovieRequest struct {
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}
