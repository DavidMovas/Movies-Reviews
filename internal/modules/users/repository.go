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

func (r Repository) GetUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	user := UserWithPassword{}
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}
