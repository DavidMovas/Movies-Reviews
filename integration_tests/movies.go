package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/DavidMovas/Movies-Reviews/contracts"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
)

var (
	godFather *contracts.MovieDetails
	starWars  *contracts.MovieDetails
	titanic   *contracts.MovieDetails
)

func moviesAPIChecks(t *testing.T, c *client.Client, _ *config.Config) {
	t.Run("movies.GetMovies: empty: success", func(t *testing.T) {
		req := &contracts.GetMoviesRequest{}
		res, err := c.GetMovies(req)
		require.NoError(t, err)
		require.Equal(t, 0, res.Total)
	})

	t.Run("movies.GetMovieById: not found", func(t *testing.T) {
		_, err := c.GetMovieByID(1)
		requireNotFoundError(t, err, "movie", "id", 1)
	})

	t.Run("movies.CreateMovie: insufficient permissions", func(t *testing.T) {
		req := &contracts.CreateMovieRequest{}
		_, err := c.CreateMovie("", req)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("movies.CreateMovie: create 3 movies: success", func(t *testing.T) {
		cases := []struct {
			req  *contracts.CreateMovieRequest
			addr **contracts.MovieDetails
		}{
			{
				req: &contracts.CreateMovieRequest{
					Title:       "The Godfather",
					ReleaseDate: time.Date(1972, 3, 24, 0, 0, 0, 0, time.UTC),
					Description: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.",
				},
				addr: &godFather,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Star Wars: Episode IV - A New Hope",
					ReleaseDate: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
					Description: "A young Luke Skywalker has been chosen as a pilot. However, his past leads him to become a legendary pilot of the Star Wars warship, the Imperial Starfleet.",
				},
				addr: &starWars,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Titanic",
					ReleaseDate: time.Date(1997, 12, 19, 0, 0, 0, 0, time.UTC),
					Description: "A seventeen-year-old aristocrat falls in love with a kind but poor artist aboard the luxurious, ill-fated R.M.S. Titanic.",
				},
				addr: &titanic,
			},
		}

		for _, cc := range cases {
			movie, err := c.CreateMovie(johnMooreToken, cc.req)
			require.NoError(t, err)

			*cc.addr = movie
			require.NotEmpty(t, movie.ID)
			require.NotEmpty(t, movie.CreatedAt)
		}
	})

	t.Run("movies.CreateMovie: title empty", func(t *testing.T) {
		req := &contracts.CreateMovieRequest{
			ReleaseDate: time.Date(1972, 3, 24, 0, 0, 0, 0, time.UTC),
			Description: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.",
		}
		_, err := c.CreateMovie(johnMooreToken, req)
		requireBadRequestError(t, err, "Title: less than min")
	})

	t.Run("movies.GetMovieById: success", func(t *testing.T) {
		movie, err := c.GetMovieByID(godFather.ID)
		require.NoError(t, err)
		require.Equal(t, godFather, movie)
	})

	t.Run("movies.GetMovies: success", func(t *testing.T) {
		req := &contracts.GetMoviesRequest{}
		res, err := c.GetMovies(req)

		require.NoError(t, err)
		require.Equal(t, 3, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationDefaultSize, res.Size)
		require.Equal(t, len([]*contracts.MovieDetails{godFather, starWars}), len(res.Items))

		req.Page = res.Page + 1
		res, err = c.GetMovies(req)
		require.NoError(t, err)
		require.Equal(t, 3, res.Total)
		require.Equal(t, 2, res.Page)
		require.Equal(t, testPaginationDefaultSize, res.Size)
		require.Equal(t, len([]*contracts.MovieDetails{titanic}), len(res.Items))
	})

	t.Run("movies.UpdateMovie: insufficient permissions", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			Title: ptr("The Godfather 2"),
		}
		_, err := c.UpdateMovieByID("", req, 1)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("movies.UpdateMovie: not found", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			Title: ptr("The Godfather 2"),
		}
		_, err := c.UpdateMovieByID(johnMooreToken, req, 100)
		requireNotFoundError(t, err, "movie", "id", 100)
	})

	t.Run("movies.UpdateMovie: success", func(t *testing.T) {
		releaseTime := ptr(time.Date(1975, 3, 24, 0, 0, 0, 0, time.UTC))
		req := &contracts.UpdateMovieRequest{
			Title:       ptr("The Godfather 2"),
			ReleaseDate: releaseTime,
		}
		movie, err := c.UpdateMovieByID(johnMooreToken, req, godFather.ID)
		require.NoError(t, err)
		require.Equal(t, "The Godfather 2", movie.Title)
		require.Equal(t, *releaseTime, movie.ReleaseDate)
	})

	t.Run("movies.DeleteMovie: insufficient permissions", func(t *testing.T) {
		err := c.DeleteMovieByID("", 1)
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("movies.DeleteMovie: not found", func(t *testing.T) {
		err := c.DeleteMovieByID(johnMooreToken, 100)
		requireNotFoundError(t, err, "movie", "id", 100)
	})

	t.Run("movies.DeleteMovie: success", func(t *testing.T) {
		err := c.DeleteMovieByID(johnMooreToken, godFather.ID)
		require.NoError(t, err)
	})
}
