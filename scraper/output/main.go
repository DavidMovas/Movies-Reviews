package main

import (
	"log/slog"
	"os"

	"github.com/DavidMovas/Movies-Reviews/scraper/cmd"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	scraper := cmd.NewScrapCmd(logger)
	ingester := cmd.NewIngestCmd(logger)

	if err := scraper.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	if err := ingester.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
