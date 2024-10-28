package movies

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/stars"

	"github.com/DavidMovas/Movies-Reviews/contracts"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	"github.com/DavidMovas/Movies-Reviews/internal/pagination"

	"github.com/labstack/echo/v4"
)

const (
	paramMovieID   = "movieId"
	invalidMovieID = "invalid movie id"
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

// GetMovies godoc
// @Summary      Get movies
// @Description  Get movies
// @ID           get-movies
// @Tags         movies
// @Produce      json
// @Param        request body contracts.PaginatedRequestOrdered false "Request, if request body empty, default values will be used"
// @Success      200 {object} pagination.PaginatedResponseOrdered[contracts.Movie] "PaginatedResponse of Movies, total number of movies, or nil if none found"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies [get]
func (h *Handler) GetMovies(c echo.Context) error {
	req, err := echox.BindAndValidate[GetMoviesRequest](c)
	if err != nil {
		return err
	}

	pagination.SetDefaultsOrdered(&req.PaginatedRequestOrdered, h.paginationConfig)
	offset, limit := pagination.OffsetLimit(&req.PaginatedRequest)
	if err = contracts.ValidateSortRequest(req.Sort); err != nil {
		req.Sort = "id"
	}

	movies, total, err := h.service.GetMovies(c.Request().Context(), offset, limit, req.Sort, req.Order)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pagination.ResponseOrdered[*Movie](&req.PaginatedRequestOrdered, total, movies))
}

// GetMovieByID godoc
// @Summary      Get movie by id
// @Description  Get movie by id
// @ID           get-movie-by-id
// @Tags         movies
// @Produce      json
// @Param        movieId path int true "Movie ID"
// @Success      200 {object} contracts.MovieDetails "Movie details"
// @Failure      400 {object} apperrors.Error "Invalid movie id, invalid parameter or missing parameter"
// @Failure      404 {object} apperrors.Error "Movie not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies/{movieId} [get]
func (h *Handler) GetMovieByID(c echo.Context) error {
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	movie, err := h.service.GetMovieByID(c.Request().Context(), movieID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

// GetStarsByMovieID godoc
// @Summary      Get stars by movie id
// @Description  Get stars by movie id
// @ID           get-stars-by-movie-id
// @Tags         movies
// @Produce      json
// @Param        movieId path int true "Movie ID"
// @Success      200 {object} []contracts.Star "Stars for movie"
// @Failure      400 {object} apperrors.Error "Invalid movie id, invalid parameter or missing parameter"
// @Failure      404 {object} apperrors.Error "Movie not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies/{movieId}/stars [get]
func (h *Handler) GetStarsByMovieID(c echo.Context) error {
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	associatedStars, err := h.service.GetStarsByMovieID(c.Request().Context(), movieID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, associatedStars)
}

// CreateMovie godoc
// @Summary      Create movie
// @Description  Create movie
// @ID           create-movie
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        request body contracts.CreateMovieRequest true "Movie details"
// @Success      201 {object} contracts.MovieDetails "Movie created"
// @Failure      403 {object} apperrors.Error "Insufficient permissions"
// @Failure      400 {object} apperrors.Error "Invalid request, invalid parameter or missing parameter"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies [post]
func (h *Handler) CreateMovie(c echo.Context) error {
	req, err := echox.BindAndValidate[CreateMovieRequest](c)
	if err != nil {
		return err
	}

	movie := &MovieDetails{
		Movie: Movie{
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
		},
		Description: req.Description,
	}

	for _, genreID := range req.GenreIDs {
		movie.Genres = append(movie.Genres, &genres.Genre{
			ID: genreID,
		})
	}

	for _, creditID := range req.Cast {
		movie.Cast = append(movie.Cast, &MovieCredit{
			Star: stars.Star{
				ID: creditID.StarID,
			},
			Role:    creditID.Role,
			Details: creditID.Details,
		})
	}

	movie, err = h.service.CreateMovie(c.Request().Context(), movie)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, movie)
}

// UpdateMovieByID godoc
// @Summary      Update movie by id
// @Description  Update movie by id
// @ID           update-movie-by-id
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param        movieId path int true "Movie ID"
// @Param        request body contracts.UpdateMovieRequest true "Updated movie details, fields that are not provided will not be updated"
// @Success      200 {object} contracts.MovieDetails "Movie updated"
// @Failure      400 {object} apperrors.Error "Invalid movie id, invalid parameter or missing parameter"
// @Failure      403 {object} apperrors.Error "Insufficient permissions"
// @Failure      404 {object} apperrors.Error "Movie not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies/{movieId} [put]
func (h *Handler) UpdateMovieByID(c echo.Context) error {
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	req, err := echox.BindAndValidate[UpdateMovieRequest](c)
	if err != nil {
		return err
	}

	movie, err := h.service.UpdateMovieByID(c.Request().Context(), movieID, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, movie)
}

// DeleteMovieByID godoc
// @Summary      Delete movie by id
// @Description  Delete movie by id
// @ID           delete-movie-by-id
// @Tags         movies
// @Produce      json
// @Param        movieId path int true "Movie ID"
// @Success      200 "Movie deleted (softly deleting)"
// @Failure      400 {object} apperrors.Error "Invalid movie id, invalid parameter or missing parameter"
// @Failure      403 {object} apperrors.Error "Insufficient permissions"
// @Failure      404 {object} apperrors.Error "Movie not found"
// @Failure      500 {object} apperrors.Error "Internal server error"
// @Router       /movies/{movieId} [delete]
func (h *Handler) DeleteMovieByID(c echo.Context) error {
	movieID, err := echox.ReadFromParam[int](c, paramMovieID, invalidMovieID)
	if err != nil {
		return err
	}

	if err = h.service.DeleteMovieByID(c.Request().Context(), movieID); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
