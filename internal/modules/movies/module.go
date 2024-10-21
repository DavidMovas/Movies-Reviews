package movies

import "github.com/jackc/pgx/v4/pgxpool"

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(db *pgxpool.Pool) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}
