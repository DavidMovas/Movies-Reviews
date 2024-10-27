package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/server"

	_ "github.com/DavidMovas/Movies-Reviews/docs"
)

var dbGracefulTime = time.Second * 10

// @title Movies Reviews API
// @version 1.0.0

// @description This API is used for managing movies reviews.
// @termsOfService http://swagger.io/terms/

// @contact.name API Creator Support
// @contact.url https://github.com/DavidMovas/Movies-Reviews
// @contact.email david417program@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http
func main() {
	cfg, err := config.NewConfig()
	failOnError(err, "failed to load config")

	srv, err := server.New(context.Background(), cfg)
	failOnError(err, "failed to start server")

	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

		<-signalCh
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), dbGracefulTime)
		defer cancel()

		if err = srv.Shutdown(ctx); err != nil {
			slog.Warn("Server forced to shutdown", "error", err)
		} else {
			slog.Info("Server shutdown")
		}
	}()

	if err = srv.Start(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped")
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error("Error", err.Error(), msg)
		os.Exit(1)
	}
}
