package client

import "github.com/DavidMovas/Movies-Reviews/contracts"

func (c *Client) GetGenres() ([]*contracts.Genre, error) {
	var genres []*contracts.Genre

	_, err := c.client.R().
		SetResult(&genres).
		Get(c.path("/api/genres"))

	return genres, err
}

func (c *Client) GetGenreById(genreId int) (*contracts.Genre, error) {
	var genre *contracts.Genre

	_, err := c.client.R().
		SetResult(&genre).
		Get(c.path("/api/genres/%d", genreId))

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

func (c *Client) UpdateGenreById(accessToken string, req *contracts.UpdateGenreRequest, genreId int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Put(c.path("/api/genres/%d", genreId))

	return err
}

func (c *Client) DeleteGenreById(accessToken string, genreId int) error {
	_, err := c.client.R().
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		Delete(c.path("/api/genres/%d", genreId))

	return err
}
