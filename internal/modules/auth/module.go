package auth

import (
	"github.com/DavidMovas/Movies-Reviews/internal/jwt"
	"github.com/DavidMovas/Movies-Reviews/internal/modules/users"
)

type Module struct {
	Handler    *Handler
	Service    *Service
	Repository *Repository
}

func NewModule(jwtService *jwt.Service, userService *users.Service) *Module {
	repo := NewRepository()
	service := NewService(userService, jwtService)
	handler := NewHandler(service, userService)

	return &Module{
		Handler:    handler,
		Service:    service,
		Repository: repo,
	}
}
