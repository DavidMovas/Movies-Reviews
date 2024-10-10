package users

import "time"

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type UserWithPassword struct {
	*User
	PasswordHash string
}

func newUserWithPassword() *UserWithPassword {
	return &UserWithPassword{
		User: &User{},
	}
}
