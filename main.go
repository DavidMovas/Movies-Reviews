package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
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
	e := echo.New()

	validation.SetupValidators()

	cfg, err := config.NewConfig()
	failOnError(err, "failed to load config")

	db, err := getDB(context.Background(), cfg.DBUrl)
	failOnError(err, "failed to connect to db")

	usersModule := users.NewModule(db)

	accessTime, err := time.ParseDuration(cfg.JWT.AccessExpiration)
	if err != nil {
		accessTime = time.Duration(5) * time.Minute
	}
	jwtService := jwt.NewService(cfg.JWT.Secret, accessTime)

	authModule := auth.NewModule(usersModule.Service, jwtService)
	authMiddleware := jwt.NewAuthMiddleware(cfg.JWT.Secret)

	apiGroup := e.Group("/api")

	//ENDPOINTS: auth
	apiGroup.POST("/auth/register", authModule.Handler.Register)
	apiGroup.POST("/auth/login", authModule.Handler.Login)

	//ENDPOINTS: users
	apiGroup.GET("/users", usersModule.Handler.GetExistingUsers)
	apiGroup.GET("/users/:userId", usersModule.Handler.GetExistingUserById)
	apiGroup.GET("/users/username/:username", usersModule.Handler.GetExistingUserByUsername)
	apiGroup.PUT("/users/:userId", usersModule.Handler.UpdateExistingUserById, authMiddleware, auth.Self)
	apiGroup.PUT("/users/:userId/role/:role", usersModule.Handler.UpdateUserRoleById, authMiddleware, auth.Admin)
	apiGroup.DELETE("/users/:userId", usersModule.Handler.DeleteExistingUserById, authMiddleware, auth.Self)

	go func() {
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)

		<-signalCh

		e.Logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), dbGracefulTime)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Errorf("Server forced to shutdown: %v", err)
		} else {
			e.Logger.Info("Server gracefully stopped")
		}
	}()

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		e.Logger.Fatal(err)
	}

	log.Printf("Server closed")
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
		log.Fatalf("%s: %s", msg, err)
	}
}
