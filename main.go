package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	cfg, err := config.NewConfig()
	failOnError(err, "failed to load config")

	//TODO: Add signal's listener in goroutine to shutdown gracefully
	go func() {
		signalCh := make(chan os.Signal)
		signal.Notify(signalCh, os.Interrupt, os.Kill)

		ctx, cancel := context.WithTimeout(context.Background(), 10)
		ctx, cancel = signal.NotifyContext(ctx, os.Interrupt, os.Kill)
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

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
