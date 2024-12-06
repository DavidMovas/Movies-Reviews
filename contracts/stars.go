package contracts

import "time"

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

type StarV2 struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"firstName"`
	MiddleName *string    `json:"middleName,omitempty"`
	LastName   string     `json:"lastName"`
	AvatarURL  *string    `json:"avatarUrl,omitempty"`
	CreatedAt  time.Time  `json:"createdAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

type GetStarRequest struct {
	StarID int `json:"" param:"starId" validate:"nonzero"`
}

type GetStarsRequest struct {
	PaginatedRequestOrdered
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
