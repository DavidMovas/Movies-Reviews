package echox

import (
	"fmt"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
)

type LoginRequest struct {
	Email    *string `json:"email"`
	Username *string `json:"username"`
	Password string  `json:"password"`
}

type EmailField struct {
	Email *string `json:"email" validate:"email"`
}

type UsernameField struct {
	Username *string `json:"username" validate:"min=3,max=24"`
}

type PasswordField struct {
	Password string `json:"password" validate:"password"`
}

func BindAndValidate[T any](c echo.Context) (*T, error) {
	req := new(T)

	if err := c.Bind(req); err != nil {
		return nil, apperrors.BadRequestHidden(err, "invalid or malformed request")
	}

	if err := validator.Validate(req); err != nil {
		return nil, apperrors.BadRequest(err)
	}

	return req, nil
}

func BindAndValidateLoginRequest(c echo.Context) (*LoginRequest, error) {
	var request LoginRequest

	if err := c.Bind(&request); err != nil {
		return nil, apperrors.BadRequestHidden(err, "invalid or malformed request")
	}

	if request.Email == nil && request.Username == nil {
		return nil, apperrors.BadRequest(fmt.Errorf("email or username must be provided"))
	}

	if err := validator.Validate(&PasswordField{Password: request.Password}); err != nil {
		return nil, apperrors.BadRequest(err)
	}

	if request.Email != nil {
		if err := validator.Validate(&EmailField{Email: request.Email}); err != nil {
			return nil, apperrors.BadRequest(err)
		}

		return &request, nil
	}

	if request.Username != nil {
		if err := validator.Validate(&UsernameField{Username: request.Username}); err != nil {
			return nil, apperrors.BadRequest(err)
		}

		return &request, nil
	}

	return &request, nil
}
