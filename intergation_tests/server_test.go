package intergation_tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/server"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	prepareInfrastructure(t, runServer)
}

func runServer(t *testing.T, pgConnString string) {
	cfg := &config.Config{
		DBUrl: pgConnString,
		Port:  0,
		JWT: config.JWTConfig{
			Secret:           "secret",
			AccessExpiration: time.Minute * 10,
		},
		Admin: config.AdminConfig{
			Username: "admin",
			Email:    "admin@mail.com",
			Password: "admin",
		},
		Local: true,
		Logger: config.LoggerConfig{
			Level: "info",
		},
	}

	srv, err := server.New(context.Background(), cfg)
	require.NoError(t, err)
	defer srv.Close()

	go func() {
		if serr := srv.Start(); !errors.Is(serr, http.ErrServerClosed) {
			require.NoError(t, serr)
		}
	}()

	var port int
	retry.Run(t, func(r *retry.R) {
		port, err = srv.Port()
		if err != nil {
			require.NoError(r, err)
		}
	})

	tests(t, port, cfg)

	err = srv.Shutdown(context.Background())
	require.NoError(t, err)
}

func tests(t *testing.T, port int, cfg *config.Config) {
	addr := fmt.Sprintf("http://localhost:%d", port)
	c := client.New(addr)

	t.Run("users.GetExistingUserByUsername: not found", func(t *testing.T) {
		_, err := c.GetUserByUsername("someTestName")
		requireNotFoundError(t, err, "user", "username", "someTestName")
	})

	t.Run("users.GetExistingUserById: not found", func(t *testing.T) {
		_, err := c.GetUserById(100)
		requireNotFoundError(t, err, "user", "id", 100)
	})

	var (
		johnMoore     *contracts.User
		johnMoorePass = standardPassword
	)

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

	var johnMooreToken string
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

	t.Run("users.UpdateExistingUserById: success", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID,
			Username: fmt.Sprintf("%sTEST", johnMoore.Username),
			Password: fmt.Sprintf("%sTEST", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		require.NoError(t, err)
	})

	t.Run("users.UpdateExistingUserById: another user", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID + 100,
			Username: fmt.Sprintf("%sTEST", johnMoore.Username),
			Password: fmt.Sprintf("%sTEST", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, johnMooreToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.UpdateExistingUserById: non-authenticated", func(t *testing.T) {
		req := &contracts.UpdateUserRequest{
			UserId:   johnMoore.ID + 1,
			Username: fmt.Sprintf("%sTEST", johnMoore.Username),
			Password: fmt.Sprintf("%sTEST", johnMoorePass),
		}
		err := c.UpdateUserData(contracts.NewAuthenticated(req, ""))
		requireForbiddenError(t, err, "insufficient permissions")
	})
}

var (
	standardPassword = "sgwva3!ekfRRR"
)

func requireNotFoundError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.NotFound(subject, key, value).Error()
	requireApiError(t, err, http.StatusNotFound, msg)
}

func requireUnauthorizedError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusUnauthorized, msg)
}

func requireForbiddenError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusForbidden, msg)
}

func requireBadRequestError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusBadRequest, msg)
}

func requireApiError(t *testing.T, err error, statusCode int, msg string) {
	var cerr *client.Error
	ok := errors.As(err, &cerr)
	require.True(t, ok, "expected client.Error")
	require.Equal(t, statusCode, cerr.Code)
	require.Contains(t, cerr.Message, msg)
}
