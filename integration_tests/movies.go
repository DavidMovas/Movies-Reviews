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
		req := &contracts.GetMovieRequest{MovieID: 1000}
		_, err := c.GetMovieByID(req)
		requireNotFoundError(t, err, "movie", "id", 1000)
	})

	t.Run("movies.CreateMovie: insufficient permissions", func(t *testing.T) {
		req := &contracts.CreateMovieRequest{}
		_, err := c.CreateMovie(contracts.NewAuthenticated(req, ""))
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
					GenreIDs:    []int{dramaGenre.ID, actionGenre.ID},
					Cast: []contracts.MovieCreditInfo{
						{
							StarID:  jackStar.ID,
							Role:    "director",
							Details: "Director of the movie",
						},
						{
							StarID:  jackStar.ID,
							Role:    "actor",
							Details: "Actor of the movie",
						},
						{
							StarID:  denzelStar.ID,
							Role:    "actor",
							Details: "Actor of the movie",
						},
					},
				},
				addr: &godFather,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Star Wars: Episode IV - A New Hope",
					ReleaseDate: time.Date(1977, 5, 25, 0, 0, 0, 0, time.UTC),
					Description: "A young Luke Skywalker has been chosen as a pilot. However, his past leads him to become a legendary pilot of the Star Wars warship, the Imperial Starfleet.",
					GenreIDs:    []int{actionGenre.ID, dramaGenre.ID},
					Cast: []contracts.MovieCreditInfo{
						{
							StarID:  denzelStar.ID,
							Role:    "director",
							Details: "Director of the movie",
						},
						{
							StarID:  sophiaStar.ID,
							Role:    "actor",
							Details: "Actor of the movie",
						},
					},
				},
				addr: &starWars,
			},
			{
				req: &contracts.CreateMovieRequest{
					Title:       "Titanic",
					ReleaseDate: time.Date(1997, 12, 19, 0, 0, 0, 0, time.UTC),
					Description: "A seventeen-year-old aristocrat falls in love with a kind but poor artist aboard the luxurious, ill-fated R.M.S. Titanic.",
					Cast: []contracts.MovieCreditInfo{
						{
							StarID:  denzelStar.ID,
							Role:    "voice actor",
							Details: "Voice actor of the movie",
						},
					},
				},
				addr: &titanic,
			},
		}

		for _, cc := range cases {
			movie, err := c.CreateMovie(contracts.NewAuthenticated(cc.req, johnMooreToken))
			require.NoError(t, err)

			*cc.addr = movie
			require.NotEmpty(t, movie.ID)
			require.Equal(t, cc.req.Title, movie.Title)
			require.Equal(t, len(cc.req.GenreIDs), len(movie.Genres))
			require.Equal(t, len(cc.req.Cast), len(movie.Cast))
			require.Equal(t, cc.req.Cast[0].StarID, movie.Cast[0].Star.ID)
		}
	})

	t.Run("movies.CreateMovie: title empty", func(t *testing.T) {
		req := &contracts.CreateMovieRequest{
			ReleaseDate: time.Date(1972, 3, 24, 0, 0, 0, 0, time.UTC),
			Description: "The aging patriarch of an organized crime dynasty transfers control of his clandestine empire to his reluctant son.",
		}
		_, err := c.CreateMovie(contracts.NewAuthenticated(req, johnMooreToken))
		requireBadRequestError(t, err, "Title: less than min")
	})

	t.Run("movies.GetMovieById: success", func(t *testing.T) {
		req := &contracts.GetMovieRequest{MovieID: godFather.ID}
		movie, err := c.GetMovieByID(req)
		require.NoError(t, err)

		deepMovieCompare(t, godFather, movie)
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

	t.Run("movies.GetMovies: text-search: not found", func(t *testing.T) {
		req := &contracts.GetMoviesRequest{
			SearchTerm: ptr("MATCHES_NOTHING"),
		}
		res, err := c.GetMovies(req)
		require.NoError(t, err)
		require.Equal(t, 0, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, testPaginationDefaultSize, res.Size)
		require.Equal(t, 0, len(res.Items))
	})

	t.Run("movies.GetMovies: title-text-search: success", func(t *testing.T) {
		req := &contracts.GetMoviesRequest{
			SearchTerm: ptr("Godfather"),
		}
		res, err := c.GetMovies(req)
		require.NoError(t, err)
		require.Equal(t, 1, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, godFather.Title, res.Items[0].Title)
		require.Equal(t, testPaginationDefaultSize, res.Size)
		require.Equal(t, 1, len(res.Items))
	})

	t.Run("movies.GetMovies: description-text-search: success", func(t *testing.T) {
		req := &contracts.GetMoviesRequest{
			SearchTerm: ptr("aristocrat & love"),
		}
		res, err := c.GetMovies(req)
		require.NoError(t, err)
		require.Equal(t, 1, res.Total)
		require.Equal(t, 1, res.Page)
		require.Equal(t, titanic.Title, res.Items[0].Title)
		require.Equal(t, testPaginationDefaultSize, res.Size)
		require.Equal(t, 1, len(res.Items))
	})

	t.Run("movies.GetStarsByMovieId: movie not found", func(t *testing.T) {
		req := &contracts.GetMovieRequest{MovieID: 100}
		_, err := c.GetStarsByMovieID(req)
		requireNotFoundError(t, err, "movie", "id", 100)
	})

	t.Run("movies.GetStarsByMovieId: success", func(t *testing.T) {
		req := &contracts.GetMovieRequest{MovieID: godFather.ID}
		stars, err := c.GetStarsByMovieID(req)
		require.NoError(t, err)
		require.Equal(t, 3, len(stars))
	})

	t.Run("movies.UpdateMovie: insufficient permissions", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			MovieID: 1,
			Title:   ptr("The Godfather 2"),
		}
		_, err := c.UpdateMovieByID(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("movies.UpdateMovie: not found", func(t *testing.T) {
		req := &contracts.UpdateMovieRequest{
			MovieID: 100,
			Title:   ptr("The Godfather 2"),
		}
		_, err := c.UpdateMovieByID(contracts.NewAuthenticated(req, johnMooreToken))
		requireNotFoundError(t, err, "movie", "id", req.MovieID)
	})

	t.Run("movies.UpdateMovie: success", func(t *testing.T) {
		releaseTime := ptr(time.Date(1975, 3, 24, 0, 0, 0, 0, time.UTC))
		req := &contracts.UpdateMovieRequest{
			MovieID:     godFather.ID,
			Title:       ptr("The Godfather 2"),
			ReleaseDate: releaseTime,
			GenreIDs:    []*int{&actionGenre.ID, &dramaGenre.ID, &comedyGenre.ID},
			Cast: []*contracts.MovieCreditInfo{
				{
					StarID:  denzelStar.ID,
					Role:    "director",
					Details: "Director of the movie",
				},
				{
					StarID:  sophiaStar.ID,
					Role:    "actor",
					Details: "actor of the movie",
				},
			},
		}
		movie, err := c.UpdateMovieByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
		require.Equal(t, "The Godfather 2", movie.Title)
		require.Equal(t, *releaseTime, movie.ReleaseDate)
		var movieGenreName []string
		for _, genre := range movie.Genres {
			movieGenreName = append(movieGenreName, genre.Name)
		}
		require.Equal(t, []string{actionGenre.Name, dramaGenre.Name, comedyGenre.Name}, movieGenreName)
		for i, cast := range movie.Cast {
			require.Equal(t, req.Cast[i].StarID, cast.Star.ID)
		}
	})

	t.Run("movies.DeleteMovie: insufficient permissions", func(t *testing.T) {
		req := &contracts.DeleteMovieRequest{
			MovieID: 1,
		}
		err := c.DeleteMovieByID(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("movies.DeleteMovie: not found", func(t *testing.T) {
		req := &contracts.DeleteMovieRequest{
			MovieID: 100,
		}
		err := c.DeleteMovieByID(contracts.NewAuthenticated(req, johnMooreToken))
		requireNotFoundError(t, err, "movie", "id", req.MovieID)
	})

	t.Run("movies.DeleteMovie: success", func(t *testing.T) {
		req := &contracts.DeleteMovieRequest{
			MovieID: godFather.ID,
		}
		err := c.DeleteMovieByID(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})
}

func deepMovieCompare(t *testing.T, expected, actual *contracts.MovieDetails) {
	require.Equal(t, expected.Title, actual.Title)
	require.Equal(t, expected.ReleaseDate, actual.ReleaseDate)
	require.Equal(t, expected.Description, actual.Description)
	for i, genre := range expected.Genres {
		require.Equal(t, genre.Name, actual.Genres[i].Name)
	}
	for i, cast := range expected.Cast {
		require.Equal(t, cast.Star.ID, actual.Cast[i].Star.ID)
	}
}
