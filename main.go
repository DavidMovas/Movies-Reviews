package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/log"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/auth"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/DavidMovas/Movies-Reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
)

var (
	dbConnectionTime = time.Second * 10
	dbGracefulTime   = time.Second * 10
)

func main() {
	validation.SetupValidators()

	cfg, err := config.NewConfig()
	failOnError(err, "failed to load config")

	logger, err := log.SetupLogger(cfg.Local, cfg.Logger.Level)
	failOnError(err, "failed to setup logger")
	slog.SetDefault(logger)

	db, err := getDB(context.Background(), cfg.DBUrl)
	failOnError(err, "failed to connect to db")

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler
	e.HideBanner = true
	e.HidePort = true

	usersModule := users.NewModule(db)

	accessTime, err := time.ParseDuration(cfg.JWT.AccessExpiration)
	if err != nil {
		accessTime = time.Duration(5) * time.Minute
	}

	jwtService := jwt.NewService(cfg.JWT.Secret, accessTime)

	//TODO: Create admin user from config for server starting

	authModule := auth.NewModule(usersModule.Service, jwtService)
	apiGroup := e.Group("/api")

	apiGroup.Use(jwt.NewAuthMiddleware(cfg.JWT.Secret))

	//ENDPOINTS: auth
	apiGroup.POST("/auth/register", authModule.Handler.Register)
	apiGroup.POST("/auth/login", authModule.Handler.Login)

	//ENDPOINTS: users
	apiGroup.GET("/users", usersModule.Handler.GetExistingUsers)
	apiGroup.GET("/users/:userId", usersModule.Handler.GetExistingUserById)
	apiGroup.GET("/users/username/:username", usersModule.Handler.GetExistingUserByUsername)
	apiGroup.PUT("/users/:userId", usersModule.Handler.UpdateExistingUserById, auth.Self)
	apiGroup.PUT("/users/:userId/role/:role", usersModule.Handler.UpdateUserRoleById, auth.Admin)
	apiGroup.DELETE("/users/:userId", usersModule.Handler.DeleteExistingUserById, auth.Self)

	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-signalCh
		slog.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), dbGracefulTime)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			slog.Warn("Server forced to shutdown", "error", err)
		} else {
			slog.Info("Server shutdown")
		}
	}()

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}

	slog.Info("Server stopped")
}

func getDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbConnectionTime)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return db, err
}

func failOnError(err error, msg string) {
	if err != nil {
		slog.Error("Error", err, msg)
		os.Exit(1)
	}
}
