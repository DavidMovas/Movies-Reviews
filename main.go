package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo"
)

var (
	dbConnectionTime = time.Second * 10
)

func main() {
	e := echo.New()

	cfg, err := config.NewConfig()
	failOnError(err, "failed to load config")

	db, err := getDB(context.Background(), cfg.DBUrl)
	failOnError(err, "failed to connect to db")

	if err := db.Ping(context.Background()); err != nil {
		failOnError(err, "failed to ping db")
	}

	//TODO: Add signal's listener in goroutine to shutdown gracefully
	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, os.Kill)

		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()

		for {
			select {
			case <-signalCh:
				e.Logger.Info("Shutting down server...")
				if err := e.Shutdown(ctx); err != nil {
					e.Logger.Error(err)
				}
			}
		}
	}()

	if err := e.Start(fmt.Sprintf(":%d", cfg.Port)); err != nil {
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
