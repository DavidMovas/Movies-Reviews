package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func (r Repository) Create(ctx context.Context, user *UserWithPassword) (err error) {
	r.db.QueryRow(ctx,
		`INSERT INTO users (username, email, pass_hash) 
    		VALUES ($1, $2, $3) RETURNING id`, user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID)
	return err
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}
