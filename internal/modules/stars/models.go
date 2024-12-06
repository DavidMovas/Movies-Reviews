package stars

import (
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"
)

type Star struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"firstName"`
	MiddleName *string    `json:"middleName,omitempty"`
	LastName   string     `json:"lastName"`
	AvatarURL  *string    `json:"avatarUrl,omitempty"`
	BirthDate  time.Time  `json:"birthDate"`
	BirthPlace *string    `json:"birthPlace,omitempty"`
	DeathDate  *time.Time `json:"deathDate,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
	IMDbURL    *string    `json:"imdbUrl,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

type GetStarRequest struct {
	StarID int `json:"-" param:"starId" validate:"nonzero"`
}

type GetStarsRequest struct {
	pagination.PaginatedRequest
}

type MovieStarsRelation struct {
	MovieID int
	StarID  int
	Role    string
	Details string
	OrderNo int
}

func (m MovieStarsRelation) Key() any {
	type MovieStarsRelationKey struct {
		MovieID, StarID int
		Role            string
	}

	return MovieStarsRelationKey{m.MovieID, m.StarID, m.Role}
}

type MovieCredit struct {
	Star    Star
	Role    string
	Details string
}

type CreateStarRequest struct {
	FirstName  string     `json:"firstName" validate:"min=1,max=50"`
	MiddleName *string    `json:"middleName,omitempty" validate:"max=50"`
	LastName   string     `json:"lastName" validate:"min=1,max=50"`
	AvatarURL  *string    `json:"avatarUrl,omitempty"`
	BirthDate  time.Time  `json:"birthDate" validate:"nonzero"`
	BirthPlace *string    `json:"birthPlace,omitempty" validate:"max=100"`
	DeathDate  *time.Time `json:"deathDate,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
	IMDbURL    *string    `json:"imdbUrl,omitempty"`
}

type UpdateStarRequest struct {
	StarID     int        `json:"-" param:"starId" validate:"nonzero"`
	FirstName  *string    `json:"firstName,omitempty" validate:"max=50"`
	MiddleName *string    `json:"middleName,omitempty" validate:"max=50"`
	LastName   *string    `json:"lastName,omitempty" validate:"max=50"`
	BirthDate  *time.Time `json:"birthDate,omitempty"`
	BirthPlace *string    `json:"birthPlace,omitempty" validate:"max=100"`
	DeathDate  *time.Time `json:"deathDate,omitempty"`
	Bio        *string    `json:"bio,omitempty"`
}

type DeleteStarRequest struct {
	StarID int `json:"-" param:"starId" validate:"nonzero"`
}

func NewStar() *Star {
	return &Star{}
}

func (c *CreateStarRequest) ToStar() *Star {
	return &Star{
		FirstName:  c.FirstName,
		MiddleName: c.MiddleName,
		LastName:   c.LastName,
		BirthDate:  c.BirthDate,
		BirthPlace: c.BirthPlace,
		DeathDate:  c.DeathDate,
		Bio:        c.Bio,
	}
}

func (s Star) Normalize() *Star {
	return &Star{
		ID:         s.ID,
		FirstName:  s.FirstName,
		MiddleName: normalizeString(s.MiddleName),
		LastName:   s.LastName,
		BirthDate:  s.BirthDate,
		BirthPlace: normalizeString(s.BirthPlace),
		DeathDate:  normalizeDate(s.DeathDate),
		Bio:        normalizeString(s.Bio),
		CreatedAt:  s.CreatedAt,
		DeletedAt:  normalizeDate(s.DeletedAt),
	}
}

func normalizeString(s *string) *string {
	if s == nil || *s == "" {
		return nil
	}

	return s
}

func normalizeDate(date *time.Time) *time.Time {
	if date == nil || date.IsZero() {
		return nil
	}

	return date
}
