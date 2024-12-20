package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/DavidMovas/Movies-Reviews/scraper/ingesters"

	"github.com/DavidMovas/Movies-Reviews/scraper/models"

	"github.com/DavidMovas/Movies-Reviews/client"
	"github.com/DavidMovas/Movies-Reviews/contracts"

	"github.com/spf13/cobra"
)

type IngestOptions struct {
	Input    string
	URL      string
	Email    string
	Password string
}

func NewIngestCmd(logger *slog.Logger) *cobra.Command {
	var opts IngestOptions

	cmd := &cobra.Command{
		Use:   "ingest",
		Short: "Ingest movie info",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runIngest(&opts, logger)
		},
	}

	cmd.Flags().StringVarP(&opts.Input, "input", "i", "", "Input directory")
	cmd.Flags().StringVarP(&opts.URL, "url", "u", "http://localhost:8000", "API URL")
	cmd.Flags().StringVarP(&opts.Email, "email", "e", "", "User email")
	cmd.Flags().StringVarP(&opts.Password, "password", "p", "", "User password")

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("password")

	return cmd
}

func runIngest(opts *IngestOptions, logger *slog.Logger) error {
	cl := client.New(opts.URL)
	res, err := cl.LoginUser(&contracts.LoginUserRequest{
		Email:    opts.Email,
		Password: opts.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to login ingest user: %w", err)
	}
	token := res.AccessToken
	_ = token
	logger.Info("Logged in successfully")

	var (
		genres []string
		stars  map[string]*models.Star
		bios   map[string]*models.Bio
		movies map[string]*models.Movie
		cast   map[string]*models.Cast
	)

	unmarshal := func(v any) func([]byte) error {
		return func(data []byte) error {
			return json.Unmarshal(data, v)
		}
	}

	reads := []struct {
		path      string
		unmarshal func([]byte) error
	}{
		{path: filepath.Join(opts.Input, "genres.json"), unmarshal: unmarshal(&genres)},
		{path: filepath.Join(opts.Input, "stars.json"), unmarshal: unmarshal(&stars)},
		{path: filepath.Join(opts.Input, "bios.json"), unmarshal: unmarshal(&bios)},
		{path: filepath.Join(opts.Input, "movies.json"), unmarshal: unmarshal(&movies)},
		{path: filepath.Join(opts.Input, "cast.json"), unmarshal: unmarshal(&cast)},
	}

	for _, read := range reads {
		var data []byte
		data, err = os.ReadFile(read.path)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", read.path, err)
		}
		if err = read.unmarshal(data); err != nil {
			return fmt.Errorf("failed to unmarshal %s: %w", read.path, err)
		}
	}

	logger.Info("Read data successfully")

	genreIngester := ingesters.NewGenreIngest(cl, token, logger)
	if err = genreIngester.Ingest(genres); err != nil {
		return fmt.Errorf("failed to ingest genres: %w", err)
	}
	starIngester := ingesters.NewStarIngester(cl, token, logger)
	if err = starIngester.Ingest(stars, bios); err != nil {
		return fmt.Errorf("failed to ingest stars: %w", err)
	}
	movieIngester := ingesters.NewMovieIngester(cl, token, genreIngester.Converter, starIngester.Converter, logger)
	if err = movieIngester.Ingest(movies, cast); err != nil {
		return fmt.Errorf("failed to ingest movies: %w", err)
	}

	return nil
}
