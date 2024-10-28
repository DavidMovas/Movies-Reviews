package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetGenres() ([]*contracts.Genre, error) {
	var genres []*contracts.Genre

	_, err := c.client.R().
		SetResult(&genres).
		Get(c.path("/api/genres"))

	return genres, err
}

func (c *Client) GetGenreByID(req *contracts.GetGenreRequest) (*contracts.Genre, error) {
	var genre *contracts.Genre

	_, err := c.client.R().
		SetResult(&genre).
		Get(c.path("/api/genres/%d", req.GenreID))

	return genre, err
}

func (c *Client) CreateGenre(req *contracts.AuthenticatedRequest[contracts.CreateGenreRequest]) (*contracts.Genre, error) {
	var genre *contracts.Genre

	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		SetResult(&genre).
		Post(c.path("/api/genres"))

	return genre, err
}

func (c *Client) UpdateGenreByID(req *contracts.AuthenticatedRequest[contracts.UpdateGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/genres/%d", req.Request.GenreID))

	return err
}

func (c *Client) DeleteGenreByID(req *contracts.AuthenticatedRequest[contracts.DeleteGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		Delete(c.path("/api/genres/%d", req.Request.GenreID))

	return err
}
