package contracts

import "time"

type Movie struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	ReleaseDate time.Time  `json:"releaseDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	Version     int        `json:"version"`
}

type MovieDetails struct {
	Movie
	Description string `json:"description"`
}
