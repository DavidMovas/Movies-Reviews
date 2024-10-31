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

// GetReviewsByMovieID godoc
// @Summary      Get reviews by movie ID
// @Description  Get reviews by movie ID
// @ID           get-reviews-by-movie-id
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        movieId path int true "Movie ID"
// @Param        request body contracts.GetReviewsByMovieIDRequest false "Pagination request, if request body empty, default values will be used"
// @Success      200 {object} pagination.PaginatedResponse[contracts.Review] "PaginatedResponse of Reviews, total number of reviews, or nil if none found"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies/{movieId}/reviews [get]
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

// GetReviewsByUserID godoc
// @Summary      Get reviews by user ID
// @Description  Get reviews by user ID
// @ID           get-reviews-by-user-id
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        userId path int true "User ID"
// @Param        request body contracts.GetReviewsByUserIDRequest false "Pagination request, if request body empty, default values will be used"
// @Success      200 {object} pagination.PaginatedResponse[contracts.Review] "PaginatedResponse of Reviews, total number of reviews, or nil if none found"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /users/{userId}/reviews [get]
func (h *Handler) GetReviewsByUserID(c echo.Context) error {
	req, err := echox.BindAndValidate[GetReviewsByUserIDRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaultsOrderedWith(&req.PaginatedRequestOrdered, h.paginationConfig, "created_at", "desc")
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)

	reviews, total, err := h.service.GetReviewsByUserID(c.Request().Context(), req.UserID, offset, limit, req.Sort, req.Order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.ResponseOrdered[*Review](&req.PaginatedRequestOrdered, total, reviews))
}

// GetReviewByID godoc
// @Summary      Get review by ID
// @Description  Get review by ID
// @ID           get-review-by-id
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        reviewId path int true "Review ID"
// @Success      200 {object} contracts.Review "Review"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      404 {object} apperrors.Error "Review not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /reviews/{reviewId} [get]
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

// CreateReview godoc
// @Summary      Create review
// @Description  Create review
// @ID           create-review
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        request body contracts.CreateReviewRequest true "Create review request, movieId and userId are required be unique"
// @Success      201 {object} contracts.Review "Review"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      401 {object} apperrors.Error "Unauthorized"
// @Failure      403 {object} apperrors.Error "Forbidden"
// @Failure      409 {object} apperrors.Error "Review already exists"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies/{movieId}/reviews [post]
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

// UpdateReviewByID godoc
// @Summary      Update review by ID
// @Description  Update review by ID
// @ID           update-review-by-id
// @Tags         reviews
// @Accept       json
// @Produce      json
// @Param        reviewId path int true "Review ID"
// @Param        request body contracts.UpdateReviewRequest true "Update review request, at least one field is required, if optional fields are empty, it will set default values"
// @Success      200 {object} contracts.Review "Review"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      401 {object} apperrors.Error "Unauthorized"
// @Failure      403 {object} apperrors.Error "Forbidden"
// @Failure      404 {object} apperrors.Error "Not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /reviews/{reviewId} [put]
func (h *Handler) UpdateReviewByID(c echo.Context) error {
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

// DeleteReviewByID godoc
// @Summary      Delete review by ID
// @Description  Delete review by ID
// @ID           delete-review-by-id
// @Tags         reviews
// @Accept       json
// @Param        reviewId path int true "Review ID"
// @Success      200 "Review deleted (softly deleting)"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      401 {object} apperrors.Error "Unauthorized"
// @Failure      403 {object} apperrors.Error "Forbidden"
// @Failure      404 {object} apperrors.Error "Not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /reviews/{reviewId} [delete]
func (h *Handler) DeleteReviewByID(c echo.Context) error {
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
