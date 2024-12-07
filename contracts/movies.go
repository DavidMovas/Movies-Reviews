package contracts

import (
	"errors"
	"time"
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
	Description  string         `json:"description"`
	IMDbRating   *float64       `json:"imdbRating,omitempty"`
	IMDbURL      *string        `json:"imdbUrl,omitempty"`
	Metascore    *int           `json:"metascore,omitempty"`
	MetascoreURL *string        `json:"metascoreUrl,omitempty"`
	Version      int            `json:"version"`
	Genres       []*Genre       `json:"genres"`
	Cast         []*MovieCredit `json:"cast"`
}

type MovieDetailsV2 struct {
	Movie
	Description  string           `json:"description"`
	IMDbRating   *float64         `json:"imdbRating,omitempty"`
	IMDbURL      *string          `json:"imdbUrl,omitempty"`
	Metascore    *int             `json:"metascore,omitempty"`
	MetascoreURL *string          `json:"metascoreUrl,omitempty"`
	Version      int              `json:"version"`
	Genres       []*Genre         `json:"genres"`
	Cast         []*MovieCreditV2 `json:"cast"`
}

type MovieCredit struct {
	Star     Star    `json:"star"`
	Role     string  `json:"role"`
	HeroName *string `json:"heroName,omitempty"`
	Details  string  `json:"details,omitempty"`
}

type MovieCreditV2 struct {
	Star     *StarV2 `json:"star"`
	Role     string  `json:"role"`
	HeroName *string `json:"heroName,omitempty"`
	Details  string  `json:"details,omitempty"`
}

type GetMovieRequest struct {
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}

type MovieCreditInfo struct {
	StarID   int     `json:"starId"`
	Role     string  `json:"role"`
	HeroName *string `json:"heroName,omitempty"`
	Details  string  `json:"details,omitempty"`
	IMDbURL  *string `json:"imdbUrl,omitempty"`
}

type GetMoviesRequest struct {
	PaginatedRequestOrdered
	SearchTerm *string `json:"-" query:"q"`
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

func (r *GetMoviesRequest) ToQueryParams() map[string]string {
	params := r.PaginatedRequestOrdered.ToQueryParams()
	if r.SearchTerm != nil {
		params["q"] = *r.SearchTerm
	}
	return params
}

func ValidateSortRequest(sort string) error {
	if sort != "id" && sort != "title" && sort != "releaseDate" && sort != "created_at" {
		return errors.New("invalid sort field")
	}

	return nil
}
