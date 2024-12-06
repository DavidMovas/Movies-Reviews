package starsv2

import "time"

type StarShort struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"firstName"`
	MiddleName *string    `json:"middleName,omitempty"`
	LastName   string     `json:"lastName"`
	CreatedAt  time.Time  `json:"createdAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}
