package contracts

import (
	"errors"
	"time"
)

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"releaseDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type MovieDetails struct {
	Movie
	Description string         `json:"description"`
	Version     int            `json:"version"`
	Genres      []*Genre       `json:"genres"`
	Cast        []*MovieCredit `json:"cast"`
}

type MovieCredit struct {
	Star    Star   `json:"star"`
	Role    string `json:"role"`
	Details string `json:"details,omitempty"`
}

type GetMovieRequest struct {
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}

type MovieCreditInfo struct {
	StarID  int    `json:"starId"`
	Role    string `json:"role"`
	Details string `json:"details,omitempty"`
}

type GetMoviesRequest struct {
	PaginatedRequestOrdered
	SearchTerm *string `json:"-" query:"q"`
}

type CreateMovieRequest struct {
	Title       string            `json:"title" validate:"min=1,max=100"`
	ReleaseDate time.Time         `json:"releaseDate" validate:"nonzero"`
	Description string            `json:"description"`
	GenreIDs    []int             `json:"genreIds" validate:"nonzero"`
	Cast        []MovieCreditInfo `json:"cast"`
}

type UpdateMovieRequest struct {
	MovieID     int                `json:"-" param:"movieId" validate:"nonzero"`
	Title       *string            `json:"title,omitempty" validate:"max=100"`
	ReleaseDate *time.Time         `json:"releaseDate,omitempty"`
	Description *string            `json:"description,omitempty"`
	Version     int                `json:"version"`
	GenreIDs    []*int             `json:"genreIds,omitempty"`
	Cast        []*MovieCreditInfo `json:"cast,omitempty"`
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
