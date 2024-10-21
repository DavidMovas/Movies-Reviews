package client

import (
	"github.com/DavidMovas/Movies-Reviews/contracts"
)

func (c *Client) GetMovies(req *contracts.GetMoviesRequest) (*contracts.PaginatedResponseOrdered[*contracts.Movie], error) {
	var resp *contracts.PaginatedResponseOrdered[*contracts.Movie]

	_, err := c.client.R().
		SetResult(&resp).
		SetQueryParams(req.ToQueryParams()).
		Get(c.path("/api/movies"))

	return resp, err
}

func (c *Client) GetMovieByID(movieID int) (*contracts.MovieDetails, error) {
	var movie *contracts.MovieDetails

	_, err := c.client.R().
		SetResult(&movie).
		Get(c.path("/api/movies/%d", movieID))

	return movie, err
}

func (c *Client) CreateMovie(accessToken string, req *contracts.CreateMovieRequest) (*contracts.MovieDetails, error) {
	var movie *contracts.MovieDetails

	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetResult(&movie).
		SetBody(req).
		Post(c.path("/api/movies"))

	return movie, err
}

func (c *Client) UpdateMovieByID(accessToken string, req *contracts.UpdateMovieRequest, movieID int) (*contracts.MovieDetails, error) {
	var movie *contracts.MovieDetails

	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetResult(&movie).
		SetBody(req).
		Put(c.path("/api/movies/%d", movieID))

	return movie, err
}

func (c *Client) DeleteMovieByID(accessToken string, movieID int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		Delete(c.path("/api/movies/%d", movieID))

	return err
}
