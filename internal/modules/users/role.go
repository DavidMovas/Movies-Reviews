package users

const (
	AdminRole  = "admin"
	EditorRole = "editor"
	UserRole   = "user"
)

func ValidateRole(role string) bool {
	return role == AdminRole || role == EditorRole || role == UserRole
}
