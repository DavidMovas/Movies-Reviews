package movies

import (
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"
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
	Description string          `json:"description"`
	Version     int             `json:"version"`
	Genres      []*genres.Genre `json:"genres"`
}

type GetMovieRequest struct {
	MovieID int `json:"movieId" validate:"nonzero"`
}

type GetMoviesRequest struct {
	pagination.PaginatedRequestOrdered
}

type CreateMovieRequest struct {
	Title       string    `json:"title" validate:"min=1,max=100"`
	ReleaseDate time.Time `json:"releaseDate" validate:"nonzero"`
	Description string    `json:"description"`
	GenreIDs    []int     `json:"genreIds" validate:"nonzero"`
}

type UpdateMovieRequest struct {
	MovieID     int        `json:"-" param:"movieId" validate:"nonzero"`
	Title       *string    `json:"title,omitempty" validate:"max=100"`
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	Description *string    `json:"description,omitempty"`
	Version     int        `json:"version"`
	GenreIDs    []*int     `json:"genreIds,omitempty"`
}

type DeleteMovieRequest struct {
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}
