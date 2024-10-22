package pagination

import (
	"github.com/DavidMovas/Movies-Reviews/internal/config"
)

func SetDefaults(r *PaginatedRequest, cfg *config.PaginationConfig) {
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

func SetDefaultsOrdered(r *PaginatedRequestOrdered, cfg *config.PaginationConfig) {
	SetDefaults(&r.PaginatedRequest, cfg)

	if r.Sort == "" {
		r.Sort = "id"
	}

	if r.Order == "" || r.Order != "asc" && r.Order != "desc" {
		r.Order = "asc"
	}
}

func OffsetLimit(r *PaginatedRequest) (int, int) {
	offset := (r.Page - 1) * r.Size
	limit := r.Size
	return offset, limit
}

func Response[T any](r *PaginatedRequest, total int, items []T) *PaginatedResponse[T] {
	return &PaginatedResponse[T]{
		Page:  r.Page,
		Size:  r.Size,
		Total: total,
		Items: items,
	}
}

func ResponseOrdered[T any](r *PaginatedRequestOrdered, total int, items []T) *PaginatedResponseOrdered[T] {
	return &PaginatedResponseOrdered[T]{
		PaginatedResponse: PaginatedResponse[T]{
			Page:  r.Page,
			Size:  r.Size,
			Total: total,
			Items: items,
		},
		Sort:  r.Sort,
		Order: r.Order,
	}
}
