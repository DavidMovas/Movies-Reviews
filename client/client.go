package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client  *resty.Client
	baseURL string
}

func New(url string) *Client {
	hc := &http.Client{}
	rc := resty.NewWithClient(hc)
	rc.OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
		if response.IsError() {
			herr := contracts.HTTPError{}
			_ = json.Unmarshal(response.Body(), &herr)

			return &Error{Code: response.StatusCode(), Message: herr.Message}
		}
		return nil
	})

	return &Client{
		client:  rc,
		baseURL: url,
	}
}

func (c *Client) path(f string, args ...any) string {
	return fmt.Sprintf(c.baseURL+f, args...)
}
