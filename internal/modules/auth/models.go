package auth

type RegisterUserRequest struct {
	Username string `json:"username" validate:"min=3,max=24"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"password"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token" nonzero:"true"`
}

type AuthenticatedRequest[T any] struct {
	AccessToken string
	Request     T
}

func NewAuthenticated[T any](req T, accessToken string) *AuthenticatedRequest[T] {
	return &AuthenticatedRequest[T]{
		AccessToken: accessToken,
		Request:     req,
	}
}
