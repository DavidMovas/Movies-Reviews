package stars

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/config"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"

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
	paginationConfig *config.PaginationConfig
}

func NewHandler(service *Service, paginationConfig *config.PaginationConfig) *Handler {
	return &Handler{
		service,
		paginationConfig,
	}
}

func (h *Handler) GetStars(c echo.Context) error {
	req, err := echox.BindAndValidate[contracts.GetStarsRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	stars, total, err := h.Service.GetStarsPaginated(c.Request().Context(), offset, limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.Response[*contracts.Star](&req.PaginatedRequest, total, stars))
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

	if err = h.Service.DeleteStarByID(c.Request().Context(), starID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
