package stars

import "time"

type StarV2 struct {
	ID         int        `json:"id"`
	FirstName  string     `json:"firstName"`
	MiddleName *string    `json:"middleName,omitempty"`
	LastName   string     `json:"lastName"`
	CreatedAt  time.Time  `json:"createdAt"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

func (s Star) ConvertToV2() *StarV2 {
	return &StarV2{
		ID:         s.ID,
		FirstName:  s.FirstName,
		MiddleName: normalizeString(s.MiddleName),
		LastName:   s.LastName,
		CreatedAt:  s.CreatedAt,
		DeletedAt:  normalizeDate(s.DeletedAt),
	}
}
