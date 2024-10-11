package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetExistingUsers(c echo.Context) error {
	return apperrors.BadRequest(errors.New("not implemented"))
}

func (h *Handler) GetExistingUserById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	user, err := h.service.GetExistingUserById(c.Request().Context(), userId)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) GetExistingUserByUsername(c echo.Context) error {
	username := c.Param("username")

	user, err := h.service.GetExistingUserByUsername(c.Request().Context(), username)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateExistingUserById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[NewUserData](c)
	if err != nil {
		return err
	}

	if err := h.service.UpdateExistingUserById(c.Request().Context(), userId, raq); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateUserRoleById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	newRole := c.Param("role")

	if !ValidateRole(newRole) {
		return apperrors.BadRequestHidden(errors.New("invalid role"), "role unknown")
	}

	if err := h.service.UpdateUserRoleById(c.Request().Context(), userId, newRole); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteExistingUserById(c echo.Context) error {
	userId, err := readUserId(c)
	if err != nil {
		return err
	}

	return h.service.DeleteExistingUserById(c.Request().Context(), userId)
}

func readUserId(c echo.Context) (int, error) {
	userId := c.Param("userId")
	if userId == "" {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "invalid user id")
	}

	return id, nil
}

type NewUserData struct {
	Username string `json:"username" validate:"min=3,max=24"`
	Password string `json:"password" validate:"password"`
}
