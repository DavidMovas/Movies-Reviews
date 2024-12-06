package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/reviews"

	"github.com/DavidMovas/Movies-Reviews/docs"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/echox"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/log"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/auth"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/movies"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/stars"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
	"github.com/DavidMovas/Movies-Reviews/internal/validation"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	dbConnectionTime  = time.Second * 10
	adminCreationTime = time.Second * 5
)

type Server struct {
	e       *echo.Echo
	cfg     *config.Config
	closers []func() error
}

func New(ctx context.Context, cfg *config.Config) (*Server, error) {
	logger, err := log.SetupLogger(cfg.Local, cfg.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("setuo logger: %w", err)
	}
	slog.SetDefault(logger)

	validation.SetupValidators()

	var closers []func() error
	db, err := getDB(ctx, cfg.DBUrl)
	if err != nil {
		return nil, fmt.Errorf("connecect to db: %w", err)
	}

	closers = append(closers, func() error {
		db.Close()
		return nil
	})

	jwtService := jwt.NewService(cfg.JWT.Secret, cfg.JWT.AccessExpiration)
	usersModule := users.NewModule(db)
	authModule := auth.NewModule(jwtService, usersModule.Service)
	genresModule := genres.NewModule(db)
	starsModule := stars.NewModule(db, cfg.Pagination)
	moviesModule := movies.NewModule(db, genresModule, starsModule, cfg.Pagination)
	reviewsModule := reviews.NewModule(db, moviesModule, cfg.Pagination)

	if err = createInitialAdminUser(cfg.Admin, authModule.Service); err != nil {
		return nil, withClosers(closers, fmt.Errorf("create initial admin user: %w", err))
	}

	e := echo.New()
	e.HTTPErrorHandler = echox.ErrorHandler

	e.Use(middleware.Recover())
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())

	api := e.Group("/api")
	api.Use(jwt.NewAuthMiddleware(cfg.JWT.Secret))
	api.Use(echox.Logger)

	// Swagger routes
	e.GET("/swagger*", docs.EchoSwaggerHandler)

	// Auth API routes
	api.POST("/auth/register", authModule.Handler.Register)
	api.POST("/auth/login", authModule.Handler.Login)

	// Users API routes
	api.GET("/users/:userId", usersModule.Handler.GetExistingUserByID)
	api.GET("/users/username/:username", usersModule.Handler.GetExistingUserByUsername)
	api.PUT("/users/:userId", usersModule.Handler.UpdateExistingUserByID, auth.Self)
	api.PUT("/users/:userId/role/:role", usersModule.Handler.UpdateUserRoleByID, auth.Admin)
	api.DELETE("/users/:userId", usersModule.Handler.DeleteExistingUserByID, auth.Admin)

	// Genres API routers
	api.GET("/genres", genresModule.Handler.GetGenres)
	api.GET("/genres/:genreId", genresModule.Handler.GetGenreByID)
	api.POST("/genres", genresModule.Handler.CreateGenre, auth.Editor)
	api.PUT("/genres/:genreId", genresModule.Handler.UpdateGenreByID, auth.Editor)
	api.DELETE("/genres/:genreId", genresModule.Handler.DeleteGenreByID, auth.Editor)

	// Stars API routers
	api.GET("/stars", starsModule.Handler.GetStars)
	api.GET("/stars/:starId", starsModule.Handler.GetStarByID)
	api.POST("/stars", starsModule.Handler.CreateStar, auth.Editor)
	api.PUT("/stars/:starId", starsModule.Handler.UpdateStarByID, auth.Editor)
	api.DELETE("/stars/:starId", starsModule.Handler.DeleteStarByID, auth.Editor)

	// Movies API routers
	api.GET("/movies", moviesModule.Handler.GetMovies)
	api.GET("/movies/:movieId", moviesModule.Handler.GetMovieByID)
	api.GET("/movies/v2/:movieId", moviesModule.Handler.GetMovieByIDV2)
	api.GET("/movies/:movieId/stars", moviesModule.Handler.GetStarsByMovieID)
	api.POST("/movies", moviesModule.Handler.CreateMovie, auth.Editor)
	api.PUT("/movies/:movieId", moviesModule.Handler.UpdateMovieByID, auth.Editor)
	api.DELETE("/movies/:movieId", moviesModule.Handler.DeleteMovieByID, auth.Editor)

	// Reviews API routers
	api.GET("/movies/:movieId/reviews", reviewsModule.Handler.GetReviewsByMovieID)
	api.GET("/users/:userId/reviews", reviewsModule.Handler.GetReviewsByUserID)
	api.GET("/reviews/:reviewId", reviewsModule.Handler.GetReviewByID)
	api.POST("/users/:userId/reviews", reviewsModule.Handler.CreateReview, auth.Self)
	api.PUT("/users/:userId/reviews/:reviewId", reviewsModule.Handler.UpdateReviewByID, auth.Self)
	api.DELETE("/users/:userId/reviews/:reviewId", reviewsModule.Handler.DeleteReviewByID, auth.Self)

	return &Server{
		e:       e,
		cfg:     cfg,
		closers: closers,
	}, nil
}

func (s *Server) Start() error {
	port := s.cfg.Port
	slog.Info("starting server on port", "port", port)
	return s.e.Start(fmt.Sprintf(":%d", port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func (s *Server) Close() error {
	return withClosers(s.closers, nil)
}

func (s *Server) Port() (int, error) {
	listener := s.e.Listener
	if listener == nil {
		return 0, errors.New("server is not started")
	}

	addr := listener.Addr()
	if addr == nil {
		return 0, errors.New("server is not started")
	}
	return addr.(*net.TCPAddr).Port, nil
}

func getDB(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, dbConnectionTime)
	defer cancel()

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	return db, err
}

func createInitialAdminUser(cfg config.AdminConfig, service *auth.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), adminCreationTime)
	defer cancel()

	err := service.Register(ctx, &users.User{
		Username: cfg.Username,
		Email:    cfg.Email,
		Role:     contracts.AdminRole,
	}, cfg.Password)

	switch {
	case apperrors.Is(err, apperrors.InternalCode):
		return fmt.Errorf("create initial admin user: %w", err)
	case err != nil:
		// Just ignore the error
		return nil
	default:
		slog.Info("created initial admin user", "username", cfg.Username, "email", cfg.Email)
		return nil
	}
}

func withClosers(closers []func() error, err error) error {
	errs := []error{err}

	for i := len(closers) - 1; i >= 0; i-- {
		if err = closers[i](); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
