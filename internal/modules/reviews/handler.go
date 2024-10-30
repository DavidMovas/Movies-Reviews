package reviews

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/pagination"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service          *Service
	paginationConfig *config.PaginationConfig
}

func NewHandler(service *Service, paginationConfig *config.PaginationConfig) *Handler {
	return &Handler{
		service:          service,
		paginationConfig: paginationConfig,
	}
}

func (h *Handler) GetReviewsByMovieID(c echo.Context) error {
	req, err := echox.BindAndValidate[GetReviewsByMovieIDRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaultsOrderedWith(&req.PaginatedRequestOrdered, h.paginationConfig, "created_at", "desc")
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	reviews, total, err := h.service.GetReviewsByMovieID(c.Request().Context(), req.MovieID, offset, limit, req.Sort, req.Order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.ResponseOrdered[*Review](&req.PaginatedRequestOrdered, total, reviews))
}

func (h *Handler) GetReviewsByUserID(c echo.Context) error {
	req, err := echox.BindAndValidate[GetReviewsByUserIDRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaultsOrderedWith(&req.PaginatedRequestOrdered, h.paginationConfig, "created_at", "desc")
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	reviews, total, err := h.service.GetReviewsByMovieID(c.Request().Context(), req.UserID, offset, limit, req.Sort, req.Order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.ResponseOrdered[*Review](&req.PaginatedRequestOrdered, total, reviews))
}

func (h *Handler) GetReviewByID(c echo.Context) error {
	req, err := echox.BindAndValidate[GetReviewRequest](c)
	if err != nil {
		return err
	}

	review, err := h.service.GetReviewByID(c.Request().Context(), req.ReviewID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, review)
}

func (h *Handler) CreateReview(c echo.Context) error {
	req, err := echox.BindAndValidate[CreateReviewRequest](c)
	if err != nil {
		return err
	}

	review, err := h.service.CreateReview(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, review)
}

func (h *Handler) UpdateReview(c echo.Context) error {
	req, err := echox.BindAndValidate[UpdateReviewRequest](c)
	if err != nil {
		return err
	}

	review, err := h.service.UpdateReview(c.Request().Context(), req.ReviewID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, review)
}

func (h *Handler) DeleteReview(c echo.Context) error {
	req, err := echox.BindAndValidate[DeleteReviewRequest](c)
	if err != nil {
		return err
	}

	err = h.service.DeleteReview(c.Request().Context(), req.ReviewID)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
