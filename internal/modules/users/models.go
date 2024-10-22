package users

import "time"

const (
	AdminRole  = "admin"
	EditorRole = "editor"
	UserRole   = "user"
)

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

type UpdateUserRequest struct {
	UserID   int    `json:"userId"`
	Username string `json:"username" validate:"min=3,max=24"`
	Password string `json:"password" validate:"password"`
}

func (u *User) IsDeleted() bool {
	return u.DeletedAt == nil
}

func NewUserWithPassword() *UserWithPassword {
	return &UserWithPassword{
		User: &User{},
	}
}
