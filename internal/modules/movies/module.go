package movies

import (
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/stars"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(db *pgxpool.Pool, genresModule *genres.Module, starsModule *stars.Module, paginationConfig config.PaginationConfig) *Module {
	repo := NewRepository(db, genresModule.Repository, starsModule.Repository)
	service := NewService(repo, genresModule.Repository, starsModule.Repository)
	handler := NewHandler(service, &paginationConfig)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}
