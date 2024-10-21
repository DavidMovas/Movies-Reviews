package contracts

import "time"

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"releaseDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type MovieDetails struct {
	Movie
	Description string `json:"description"`
	Version     int    `json:"version"`
}

type GetMovieRequest struct {
	MovieID int `json:"MovieId" validate:"nonzero"`
}

type GetMoviesRequest struct {
	PaginatedRequest
}

type CreateMovieRequest struct {
	Title       string    `json:"title" validate:"min=1,max=100"`
	ReleaseDate time.Time `json:"releaseDate" validate:"nonzero"`
	Description string    `json:"description"`
}

type UpdateMovieRequest struct {
	MovieID     int        `json:"-" param:"movieId" validate:"nonzero"`
	Title       *string    `json:"title,omitempty" validate:"max=100"`
	ReleaseDate *time.Time `json:"releaseDate,omitempty"`
	Description *string    `json:"description,omitempty"`
}

type DeleteMovieRequest struct {
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}
