package tests

import (
	"fmt"
	"math/rand"
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

	adminToken string
)

const (
	standardPassword = "sgwva3!ekfRRR"
)

func authAPIChecks(t *testing.T, c *client.Client, _ *config.Config) {
	t.Run("auth.RegisterUser: wrong email", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "johnmoore",
			Email:    "wrong_email",
			Password: johnMoorePass,
		}
		_, err := c.RegisterUser(req)
		requireBadRequestError(t, err, "Email: mail: missing '@' or angle-addr")
	})

	t.Run("auth.RegisterUser: wrong password", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "johnmoore",
			Email:    "johnmoore@mail.com",
			Password: "some_wrong_password",
		}
		_, err := c.RegisterUser(req)
		requireBadRequestError(t, err, "Password: password must contain at least one of the following required entries: uppercase")
	})

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

	t.Run("auth.RegisterUsers: register 5 users: success", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			registerRandomUser(t, c, "rName", "rEmail")
		}
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

	t.Run("auth.LoginUser: wrong email", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    "noexistinguser@mail.com",
			Password: standardPassword,
		}
		_, err := c.LoginUser(req)
		requireNotFoundError(t, err, "user", "email", req.Email)
	})

	t.Run("auth.LoginUser: wrong password", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    johnMoore.Email,
			Password: johnMoorePass + "wrong",
		}
		_, err := c.LoginUser(req)
		requireUnauthorizedError(t, err, "invalid password")
	})
}

func registerRandomUser(t *testing.T, c *client.Client, username, email string) *contracts.User {
	r := rand.Intn(10000)

	return register(t, c, fmt.Sprintf("%s%d", username, r), fmt.Sprintf("%s%d@mail.com", email, r), standardPassword)
}

func register(t *testing.T, c *client.Client, username, email, password string) *contracts.User {
	req := &contracts.RegisterUserRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	user, err := c.RegisterUser(req)
	require.NoError(t, err)
	return user
}

func login(t *testing.T, c *client.Client, email, password string) string {
	req := &contracts.LoginUserRequest{
		Email:    email,
		Password: password,
	}
	res, err := c.LoginUser(req)
	require.NoError(t, err)

	return res.AccessToken
}
