package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func (r Repository) Create(ctx context.Context, user *UserWithPassword) (err error) {
	err = r.db.QueryRow(ctx,
		`INSERT INTO users (username, email, pass_hash) 
    		VALUES ($1, $2, $3) RETURNING id`, user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID)

	return err
}

func (r Repository) GetExistingUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	user := newUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r Repository) GetExistingUserById(ctx context.Context, id int) (*UserWithPassword, error) {
	user := newUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role FROM users WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r Repository) DeleteExistingUserById(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	if n.RowsAffected() == 0 {
		return fmt.Errorf(" user not found")
	}

	return nil
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}
