package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetGenres() ([]*contracts.Genre, error) {
	var genres []*contracts.Genre

	_, err := c.client.R().
		SetResult(&genres).
		Get(c.path("/api/genres"))

	return genres, err
}

func (c *Client) GetGenreByID(genreID int) (*contracts.Genre, error) {
	var genre *contracts.Genre

	_, err := c.client.R().
		SetResult(&genre).
		Get(c.path("/api/genres/%d", genreID))

	return genre, err
}

func (c *Client) CreateGenre(accessToken string, req *contracts.CreateGenreRequest) (*contracts.Genre, error) {
	var genre *contracts.Genre

	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		SetResult(&genre).
		Post(c.path("/api/genres"))

	return genre, err
}

func (c *Client) UpdateGenreByID(accessToken string, req *contracts.UpdateGenreRequest, genreID int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Put(c.path("/api/genres/%d", genreID))

	return err
}

func (c *Client) DeleteGenreByID(accessToken string, genreID int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		Delete(c.path("/api/genres/%d", genreID))

	return err
}
