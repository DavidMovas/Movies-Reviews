package pagination

import (
	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/config"
)

func SetDefaults(r *contracts.PaginatedRequest, cfg *config.PaginationConfig) {
	if r.Page == 0 {
		r.Page = 1
	}

	if r.Size == 0 {
		r.Size = cfg.DefaultSize
	}

	if r.Size > cfg.MaxSize {
		r.Size = cfg.MaxSize
	}
}

func SetDefaultsOrdered(r *contracts.PaginatedRequestOrdered, cfg *config.PaginationConfig) {
	SetDefaults(&r.PaginatedRequest, cfg)

	if r.Sort == "" {
		r.Sort = "id"
	}

	if r.Order == "" || r.Order != "asc" && r.Order != "desc" {
		r.Order = "asc"
	}
}

func OffsetLimit(r *contracts.PaginatedRequest) (int, int) {
	offset := (r.Page - 1) * r.Size
	limit := r.Size
	return offset, limit
}

func Response[T any](r *contracts.PaginatedRequest, total int, items []T) *contracts.PaginatedResponse[T] {
	return &contracts.PaginatedResponse[T]{
		Page:  r.Page,
		Size:  r.Size,
		Total: total,
		Items: items,
	}
}

func ResponseOrdered[T any](r *contracts.PaginatedRequestOrdered, total int, items []T) *contracts.PaginatedResponseOrdered[T] {
	return &contracts.PaginatedResponseOrdered[T]{
		PaginatedResponse: contracts.PaginatedResponse[T]{
			Page:  r.Page,
			Size:  r.Size,
			Total: total,
			Items: items,
		},
		Sort:  r.Sort,
		Order: r.Order,
	}
}
