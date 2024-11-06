package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/DavidMovas/Movies-Reviews/scraper/collectors"
)

type ScrapOptions struct {
	Output string
}

func NewScrapCmd(logger *slog.Logger) *cobra.Command {
	var opts ScrapOptions

	cmd := &cobra.Command{
		Use:   "scrap",
		Short: "Scrap movie info",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runScrap(&opts, logger)
		},
	}

	cmd.Flags().StringVarP(&opts.Output, "output", "o", "./.", "Output directory")

	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func runScrap(opts *ScrapOptions, logger *slog.Logger) error {
	baseCollector := collectors.NewBaseCollector()

	bioCollector := collectors.NewBioCollector(collectors.Derive(baseCollector), logger)
	starCollector := collectors.NewStarCollector(collectors.Derive(baseCollector), bioCollector, logger)
	castCollector := collectors.NewCastCollector(collectors.Derive(baseCollector), starCollector, logger)
	movieCollector := collectors.NewMovieCollector(collectors.Derive(baseCollector), castCollector, logger)
	topMoviesCollector := collectors.NewTopMoviesCollector(collectors.Derive(baseCollector), movieCollector, logger)

	topMoviesCollector.Start()
	topMoviesCollector.Wait()
	movieCollector.Wait()
	castCollector.Wait()
	starCollector.Wait()
	bioCollector.Wait()

	writers := []struct {
		data any
		path string
	}{
		{data: movieCollector.Movies(), path: filepath.Join(opts.Output, "movies.json")},
		{data: movieCollector.Genres(), path: filepath.Join(opts.Output, "genres.json")},
		{data: castCollector.Cast(), path: filepath.Join(opts.Output, "cast.json")},
		{data: starCollector.Stars(), path: filepath.Join(opts.Output, "stars.json")},
		{data: bioCollector.Bios(), path: filepath.Join(opts.Output, "bios.json")},
	}

	if err := os.MkdirAll(opts.Output, os.ModePerm); err != nil {
		logger.
			With("err", err).
			Error("error creating output dir")
		return err
	}

	for _, w := range writers {
		content, err := json.Marshal(w.data)
		if err != nil {
			return fmt.Errorf("failed to marshal %s: %s", w.path, err)
		}

		if err = os.WriteFile(w.path, content, os.ModePerm); err != nil {
			return fmt.Errorf("failed to write %s: %w", w.path, err)
		}
	}

	return nil
}
