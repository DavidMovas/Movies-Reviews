package integration_tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
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

	authApiChecks(t, c, cfg)
	usersApiChecks(t, c, cfg)
}
