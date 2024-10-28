package stars

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/config"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"

	"github.com/labstack/echo/v4"
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

// GetStars godoc
// @Summary      Get stars
// @Description  Get stars
// @ID           get-stars
// @Tags         stars
// @Produce      json
// @Param        request body contracts.PaginatedRequest false "Request, if request body empty, default values will be used"
// @Success      200 {object} pagination.PaginatedResponse[contracts.Star] "PaginatedResponse of Stars, total number of stars, or nil if none found"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /stars [get]
func (h *Handler) GetStars(c echo.Context) error {
	req, err := echox.BindAndValidate[GetStarsRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaults(&req.PaginatedRequest, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	stars, total, err := h.Service.GetStarsPaginated(c.Request().Context(), offset, limit)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.Response[*Star](&req.PaginatedRequest, total, stars))
}

// GetStarByID godoc
// @Summary      Get star by id
// @Description  Get star by id
// @ID           get-star-by-id
// @Tags         stars
// @Param        starId path int true "Star ID"
// @Produce      json
// @Success      200 {object} Star "Star"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      404 {object} apperrors.Error "Star not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /stars/{starId} [get]
func (h *Handler) GetStarByID(c echo.Context) error {
	req, err := echox.BindAndValidate[GetStarRequest](c)
	if err != nil {
		return err
	}

	star, err := h.Service.GetStarByID(c.Request().Context(), req.StarID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, star)
}

// CreateStar godoc
// @Summary      Create star
// @Description  Create star
// @ID           create-star
// @Tags         stars
// @Param        request body CreateStarRequest true "Request, can have optional fields"
// @Produce      json
// @Success      201 {object} Star "Star"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      403 {object} apperrors.Error "Insufficient permissions"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /stars [post]
func (h *Handler) CreateStar(c echo.Context) error {
	raq, err := echox.BindAndValidate[CreateStarRequest](c)
	if err != nil {
		return err
	}

	star, err := h.Service.CreateStar(c.Request().Context(), raq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, star)
}

// UpdateStarByID godoc
// @Summary      Update star by id
// @Description  Update star by id
// @ID           update-star-by-id
// @Tags         stars
// @Param        starId path int true "Star ID"
// @Param        request body UpdateStarRequest true "Request, can have optional fields"
// @Produce      json
// @Success      200 {object} Star "Star"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure 	 403 {object} apperrors.Error "Insufficient permissions"
// @Failure      404 {object} apperrors.Error "Star not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /stars/{starId} [put]
func (h *Handler) UpdateStarByID(c echo.Context) error {
	raq, err := echox.BindAndValidate[UpdateStarRequest](c)
	if err != nil {
		return err
	}

	star, err := h.Service.UpdateStar(c.Request().Context(), raq.StarID, raq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, star)
}

// DeleteStarByID godoc
// @Summary      Delete star by id
// @Description  Delete star by id
// @ID           delete-star-by-id
// @Tags         stars
// @Param        starId path int true "Star ID"
// @Success      200
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      403 {object} apperrors.Error "Insufficient permissions"
// @Failure      404 {object} apperrors.Error "Star not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /stars/{starId} [delete]
func (h *Handler) DeleteStarByID(c echo.Context) error {
	req, err := echox.BindAndValidate[DeleteStarRequest](c)
	if err != nil {
		return err
	}

	if err = h.Service.DeleteStarByID(c.Request().Context(), req.StarID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
