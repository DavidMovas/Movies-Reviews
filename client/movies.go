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

func (c *Client) GetMovieByID(req *contracts.GetMovieRequest) (*contracts.MovieDetails, error) {
	var movie *contracts.MovieDetails

	_, err := c.client.R().
		SetResult(&movie).
		Get(c.path("/api/movies/%d", req.MovieID))

	return movie, err
}

func (c *Client) GetStarsByMovieID(req *contracts.GetMovieRequest) ([]*contracts.Star, error) {
	var stars []*contracts.Star

	_, err := c.client.R().
		SetResult(&stars).
		Get(c.path("/api/movies/%d/stars", req.MovieID))

	return stars, err
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

func (c *Client) UpdateMovieByID(accessToken string, req *contracts.UpdateMovieRequest) (*contracts.MovieDetails, error) {
	var movie *contracts.MovieDetails

	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetResult(&movie).
		SetBody(req).
		Put(c.path("/api/movies/%d", req.MovieID))

	return movie, err
}

func (c *Client) DeleteMovieByID(accessToken string, req *contracts.DeleteMovieRequest) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		Delete(c.path("/api/movies/%d", req.MovieID))

	return err
}
