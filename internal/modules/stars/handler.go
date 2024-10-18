package stars

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"

	"github.com/labstack/echo"
)

const (
	paramStarID   = "starId"
	invalidStarID = "invalid starId"
)

type Handler struct {
	*Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service,
	}
}

func (h *Handler) GetStars(c echo.Context) error {
	stars, err := h.Service.GetStars(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stars)
}

func (h *Handler) GetStarByID(c echo.Context) error {
	starID, err := echox.ReadFromParam[int](c, paramStarID, invalidStarID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	star, err := h.Service.GetStarByID(c.Request().Context(), starID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, star)
}

func (h *Handler) CreateStar(c echo.Context) error {
	raq, err := echox.BindAndValidate[contracts.CreateStarRequest](c)
	if err != nil {
		return err
	}

	star, err := h.Service.CreateStar(c.Request().Context(), raq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, star)
}

func (h *Handler) UpdateStarByID(c echo.Context) error {
	starID, err := echox.ReadFromParam[int](c, paramStarID, invalidStarID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[contracts.UpdateStarRequest](c)
	if err != nil {
		return err
	}

	star, err := h.Service.UpdateStar(c.Request().Context(), starID, raq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, star)
}

func (h *Handler) DeleteStarByID(c echo.Context) error {
	starID, err := echox.ReadFromParam[int](c, paramStarID, invalidStarID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	if err := h.Service.DeleteStarByID(c.Request().Context(), starID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
