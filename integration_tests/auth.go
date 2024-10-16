package integration_tests

import (
	"testing"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/stretchr/testify/require"
)

var (
	johnMoore      *contracts.User
	johnMoorePass  = standardPassword
	johnMooreToken string
)

const (
	standardPassword = "sgwva3!ekfRRR"
)

func authApiChecks(t *testing.T, c *client.Client, cfg *config.Config) {

	t.Run("auth.RegisterUser: success", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "johnmoore",
			Email:    "johnmoore@mail.com",
			Password: johnMoorePass,
		}
		user, err := c.RegisterUser(req)
		require.NoError(t, err)
		johnMoore = user

		require.Equal(t, req.Username, user.Username)
		require.Equal(t, req.Email, user.Email)
		require.Equal(t, contracts.UserRole, user.Role)
	})

	t.Run("auth.LoginUser: success", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    johnMoore.Email,
			Password: johnMoorePass,
		}
		res, err := c.LoginUser(req)
		require.NoError(t, err)
		require.NotEmpty(t, res.AccessToken)
		johnMooreToken = res.AccessToken
	})
}
