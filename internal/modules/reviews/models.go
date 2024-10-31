package reviews

import (
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"
)

type Review struct {
	ID        int        `json:"id"`
	MovieID   int        `json:"movieId"`
	UserID    int        `json:"userId"`
	Rating    int        `json:"rating"`
	Title     string     `json:"title"`
	Content   string     `json:"description"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type GetReviewRequest struct {
	ReviewID int `json:"-" param:"reviewId" validate:"nonzero"`
}

type GetReviewsByMovieIDRequest struct {
	pagination.PaginatedRequestOrdered
	MovieID int `json:"-" param:"movieId" validate:"nonzero"`
}

type GetReviewsByUserIDRequest struct {
	pagination.PaginatedRequestOrdered
	UserID int `json:"-" param:"userId" validate:"nonzero"`
}

type CreateReviewRequest struct {
	MovieID int    `json:"movieId" validate:"nonzero"`
	UserID  int    `json:"-" param:"userId" validate:"nonzero"`
	Title   string `json:"title" validate:"max=100"`
	Content string `json:"description" validate:"max=1000"`
	Rating  int    `json:"rating" validate:"min=1,max=10"`
}

type UpdateReviewRequest struct {
	UserID   int     `json:"-" param:"userId" validate:"nonzero"`
	ReviewID int     `json:"-" param:"reviewId" validate:"nonzero"`
	MovieID  int     `json:"movieId" validate:"nonzero"`
	Title    *string `json:"title,omitempty" validate:"max=100"`
	Content  *string `json:"description,omitempty" validate:"max=1000"`
	Rating   *int    `json:"rating,omitempty" validate:"min=1,max=10"`
}

type DeleteReviewRequest struct {
	UserID   int `json:"-" param:"userId" validate:"nonzero"`
	ReviewID int `json:"-" param:"reviewId" validate:"nonzero"`
}
