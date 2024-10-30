package reviews

import "github.com/DavidMovas/Movies-Reviews/internal/pagination"

type Review struct {
	ID          int     `json:"id"`
	MovieID     int     `json:"movieId"`
	UserID      string  `json:"userId"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rating      int     `json:"rating"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   *string `json:"updatedAt,omitempty"`
	DeletedAt   *string `json:"deletedAt,omitempty"`
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
	MovieID     int    `json:"-" param:"movieId" validate:"nonzero"`
	UserID      int    `json:"userId" validate:"nonzero"`
	Title       string `json:"title" validate:"max=100"`
	Description string `json:"description" validate:"max=1000"`
	Rating      int    `json:"rating" validate:"min=1,max=10"`
}

type UpdateReviewRequest struct {
	ReviewID    int     `json:"-" param:"reviewId" validate:"nonzero"`
	MovieID     int     `json:"movieId" validate:"nonzero"`
	Title       *string `json:"title,omitempty" validate:"max=100"`
	Description *string `json:"description,omitempty" validate:"max=1000"`
	Rating      *int    `json:"rating,omitempty" validate:"min=1,max=10"`
}

type DeleteReviewRequest struct {
	ReviewID int `json:"-" param:"reviewId" validate:"nonzero"`
}
