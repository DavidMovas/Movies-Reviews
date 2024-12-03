package contracts

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
	AvatarURL string     `json:"avatarUrl"`
	Bio       *string    `json:"bio,omitempty"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type UserWithPassword struct {
	*User
	PasswordHash string
}

type GetUserByIDRequest struct {
	UserID int `json:"-" param:"userId" validate:"nonzero"`
}

type GetUserByUsernameRequest struct {
	Username string `json:"-" param:"username" validate:"min=3,max=24,nonzero"`
}

type UpdateUserRoleRequest struct {
	UserID int    `json:"-" param:"userId" validate:"nonzero"`
	Role   string `json:"-" param:"role" validate:"nonzero,role"`
}

type UpdateUserRequest struct {
	UserID   int    `json:"-" param:"userId" validate:"nonzero"`
	Username string `json:"username" validate:"min=3,max=24"`
	Password string `json:"password" validate:"password"`
}

type DeleteUserRequest struct {
	UserID int `json:"-" param:"userId" validate:"nonzero"`
}

func ValidateRole(role string) bool {
	return role == AdminRole || role == EditorRole || role == UserRole
}
