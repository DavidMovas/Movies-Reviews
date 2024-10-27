package genres

import (
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/labstack/echo"

	_ "github.com/DavidMovas/Movies-Reviews/docs"
)

const (
	paramGenreID   = "genreId"
	invalidGenreID = "invalid genreId"
)

type Handler struct {
	*Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service,
	}
}

// GetGenres @Summary Get all genres
// @Description Get all genres
// @ID get-genres
// @Tags genres
// @Produce json
// @Success 200 {array} Genre "Genres, or nil if none found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /genres [get]
func (h *Handler) GetGenres(c echo.Context) error {
	genres, err := h.Service.GetGenres(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genres)
}

// GetGenreByID @Summary Get genre by id
// @Description Get genre by id
// @ID get-genre-by-id
// @Tags genres
// @Param genreId path int true "Genre ID"
// @Produce json
// @Success 200 {object} Genre "Genre"
// @Failure 400 {object} apperrors.Error "Invalid genre id, invalid parameter or missing parameter"
// @Failure 404 {object} apperrors.Error "Genre not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /genres/{genreId} [get]
func (h *Handler) GetGenreByID(c echo.Context) error {
	genreID, err := echox.ReadFromParam[int](c, paramGenreID, invalidGenreID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	genre, err := h.Service.GetGenreByID(c.Request().Context(), genreID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, genre)
}

// CreateGenre @Summary Create new genre
// @Description Create new genre
// @ID create-genre
// @Tags genres
// @Param genre body CreateGenreRequest true "Genre"
// @Produce json
// @Success 201 {object} Genre "Genre"
// @Failure 400 {object} apperrors.Error "Invalid parameter or missing parameter"
// @Failure 409 {object} apperrors.Error "Genre with that name already exists"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /genres [post]
func (h *Handler) CreateGenre(c echo.Context) error {
	raq, err := echox.BindAndValidate[CreateGenreRequest](c)
	if err != nil {
		return err
	}

	genre, err := h.Service.CreateGenre(c.Request().Context(), raq)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, genre)
}

// UpdateGenreByID @Summary Update genre by id
// @Description Update genre by id
// @ID update-genre-by-id
// @Tags genres
// @Param genreId path int true "Genre ID"
// @Param genre body UpdateGenreRequest true "Genre"
// @Produce json
// @Success 200 "Genre updated"
// @Failure 400 {object} apperrors.Error "Invalid genre id, invalid parameter or missing parameter"
// @Failure 404 {object} apperrors.Error "Genre not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /genres/{genreId} [put]
func (h *Handler) UpdateGenreByID(c echo.Context) error {
	genreID, err := echox.ReadFromParam[int](c, paramGenreID, invalidGenreID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	raq, err := echox.BindAndValidate[UpdateGenreRequest](c)
	if err != nil {
		return err
	}

	if err := h.Service.UpdateGenreByID(c.Request().Context(), genreID, raq); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// DeleteGenreByID @Summary Delete genre by id
// @Description Delete genre by id
// @ID delete-genre-by-id
// @Tags genres
// @Param genreId path int true "Genre ID"
// @Produce json
// @Success 200 "Genre deleted (softly deleting)"
// @Failure 400 {object} apperrors.Error "Invalid genre id, invalid parameter or missing parameter"
// @Failure 404 {object} apperrors.Error "Genre not found"
// @Failure 500 {object} apperrors.Error "Internal server error"
// @Router /genres/{genreId} [delete]
func (h *Handler) DeleteGenreByID(c echo.Context) error {
	genreID, err := echox.ReadFromParam[int](c, paramGenreID, invalidGenreID)
	if err != nil {
		return apperrors.BadRequest(err)
	}

	return h.Service.DeleteGenreByID(c.Request().Context(), genreID)
}
