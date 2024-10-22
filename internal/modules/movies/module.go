package movies

import (
	"github.com/DavidMovas/Movies-Reviews/internal/config"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(db *pgxpool.Pool, genresRepo *genres.Repository, paginationConfig config.PaginationConfig) *Module {
	repo := NewRepository(db, genresRepo)
	service := NewService(repo)
	handler := NewHandler(service, &paginationConfig)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}
