package integration_tests

import (
	"testing"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/stretchr/testify/require"
)

var genre *contracts.Genre

func genresApiChecks(t *testing.T, c *client.Client, _ *config.Config) {
	t.Run("genres.GetGenres: nil", func(t *testing.T) {
		genres, err := c.GetGenres()
		require.NoError(t, err)
		require.Nil(t, genres)
	})

	t.Run("genres.GetGenreById: not found", func(t *testing.T) {
		_, err := c.GetGenreById(1)
		requireNotFoundError(t, err, "genre", "id", 1)
	})

	t.Run("genres.CreateGenre: non-authenticated", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "comedy",
		}

		_, err := c.CreateGenre("", &raq)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.CreateGenre: insufficient permissions", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}

		_, err := c.CreateGenre("", &raq)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.CreateGenre: success", func(t *testing.T) {
		var err error
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}
		genre, err = c.CreateGenre(johnMooreToken, &raq)
		require.NoError(t, err)
		require.NotNil(t, genre)
	})

	t.Run("genres.CreateGenre: create 3 genres: success", func(t *testing.T) {
		raqs := []contracts.CreateGenreRequest{
			{
				Name: "comedy",
			},
			{
				Name: "drama",
			},
			{
				Name: "romance",
			},
		}

		for _, raq := range raqs {
			_, err := c.CreateGenre(johnMooreToken, &raq)
			require.NoError(t, err)
			require.NotNil(t, genre)
		}
	})

	t.Run("genres.CreateGenre: already exists", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}
		_, err := c.CreateGenre(johnMooreToken, &raq)
		requireAlreadyExistsError(t, err, "genre", "name", "action")
	})

	t.Run("genres.CreateGenre: bad request", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "a",
		}
		_, err := c.CreateGenre(adminToken, &raq)
		requireBadRequestError(t, err, "Name: less than min")
	})

	t.Run("genres.CreateGenre: insufficient permissions", func(t *testing.T) {
		raq := contracts.CreateGenreRequest{
			Name: "action",
		}
		_, err := c.CreateGenre("", &raq)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.GetGenreById: success", func(t *testing.T) {
		genre, err := c.GetGenreById(genre.ID)
		require.NoError(t, err)
		require.NotNil(t, genre)
	})

	t.Run("genres.UpdateGenreById: insufficient permissions", func(t *testing.T) {
		raq := contracts.UpdateGenreRequest{
			Name: "horror",
		}
		err := c.UpdateGenreById("", &raq, genre.ID)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.UpdateGenreById: not found", func(t *testing.T) {
		raq := contracts.UpdateGenreRequest{
			Name: "horror",
		}
		err := c.UpdateGenreById(johnMooreToken, &raq, 100)
		requireNotFoundError(t, err, "genre", "id", 100)
	})

	t.Run("genres.UpdateGenreById: already exists", func(t *testing.T) {
		raq := contracts.UpdateGenreRequest{
			Name: genre.Name,
		}
		err := c.UpdateGenreById(johnMooreToken, &raq, genre.ID+1)
		requireAlreadyExistsError(t, err, "genre", "name", genre.Name)
	})

	t.Run("genres.UpdateGenreById: success", func(t *testing.T) {
		var err error
		raq := contracts.UpdateGenreRequest{
			Name: "horror",
		}
		err = c.UpdateGenreById(johnMooreToken, &raq, genre.ID)
		require.NoError(t, err)
		genre.Name = raq.Name
	})

	t.Run("genres.DeleteGenreById: insufficient permissions", func(t *testing.T) {
		err := c.DeleteGenreById("", genre.ID)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("genres.DeleteGenreById: not found", func(t *testing.T) {
		err := c.DeleteGenreById(johnMooreToken, genre.ID+100)
		requireNotFoundError(t, err, "genre", "id", genre.ID+100)
	})

	t.Run("genres.DeleteGenreById: success", func(t *testing.T) {
		err := c.DeleteGenreById(johnMooreToken, genre.ID+3)
		require.NoError(t, err)
	})
}
