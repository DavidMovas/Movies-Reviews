package auth

import "github.com/DavidMovas/Movies-Reviews/internal/modules/users"

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(userService *users.Service) *Module {
	repo := NewRepository()
	service := NewService(userService)
	handler := NewHandler(service)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}
