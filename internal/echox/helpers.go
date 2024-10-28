package echox

import (
	"net/http"
	"strconv"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
)

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

func ReadFromParam[T any](c echo.Context, name, errMsg string) (T, error) {
	var zeroValue T

	param := c.Param(name)
	if param == "" {
		return zeroValue, echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}

	var result T
	switch any(result).(type) {
	case int:
		parsedValue, err := strconv.Atoi(param)
		if err != nil {
			return zeroValue, echo.NewHTTPError(http.StatusBadRequest, errMsg)
		}
		return any(parsedValue).(T), nil
	case float64:
		parsedValue, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return zeroValue, echo.NewHTTPError(http.StatusBadRequest, errMsg)
		}
		return any(parsedValue).(T), nil
	case string:
		return any(param).(T), nil
	default:
		return zeroValue, echo.NewHTTPError(http.StatusBadRequest, errMsg)
	}
}
