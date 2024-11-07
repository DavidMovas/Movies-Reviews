package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/DavidMovas/Movies-Reviews/scraper/cmd"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	rootCmd := &cobra.Command{
		Use: "scraper",
	}

	rootCmd.AddCommand(cmd.NewScrapCmd(logger))
	rootCmd.AddCommand(cmd.NewIngestCmd(logger))

	if err := rootCmd.Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
