package users

import (
	"errors"
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
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

func (h *Handler) GetExistingUserById(c echo.Context) error {
	userId, err := echox.ReadFromParam[int](c, "userId", "invalid userid")
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
	username, err := echox.ReadFromParam[string](c, "username", "invalid username")
	if err != nil {
		return apperrors.BadRequest(err)
	}

	user, err := h.service.GetExistingUserByUsername(c.Request().Context(), username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateExistingUserById(c echo.Context) error {
	userId, err := echox.ReadFromParam[int](c, "userId", "invalid userid")
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[contracts.UpdateUserRequest](c)
	if err != nil {
		return err
	}

	if err := h.service.UpdateExistingUserById(c.Request().Context(), userId, raq); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) UpdateUserRoleById(c echo.Context) error {
	userId, err := echox.ReadFromParam[int](c, "userId", "invalid userid")
	if err != nil {
		return apperrors.BadRequest(err)
	}

	newRole := c.Param("role")

	if !contracts.ValidateRole(newRole) {
		return apperrors.BadRequestHidden(errors.New("invalid role"), "role unknown")
	}

	if err := h.service.UpdateUserRoleById(c.Request().Context(), userId, newRole); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) DeleteExistingUserById(c echo.Context) error {
	userId, err := echox.ReadFromParam[int](c, "userId", "invalid userid")
	if err != nil {
		return err
	}

	return h.service.DeleteExistingUserById(c.Request().Context(), userId)
}
