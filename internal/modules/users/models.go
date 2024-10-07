package users

import "time"

type User struct {
	ID        uint64
	Username  string
	Email     string
	Role      string
	CreatedAt time.Time
	DeletedAt *time.Time
}

type UserWithPassword struct {
	User
	PasswordHash string
}
