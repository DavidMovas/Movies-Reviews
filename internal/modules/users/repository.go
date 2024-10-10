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

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
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
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role, created_at FROM users WHERE email = $1 AND deleted_at IS NULL`, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r Repository) GetExistingUserById(ctx context.Context, id int) (*UserWithPassword, error) {
	user := newUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role, created_at FROM users WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (r Repository) GetExistingUserByUsername(ctx context.Context, username string) (*UserWithPassword, error) {
	user := newUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role, created_at FROM users WHERE username = $1 AND deleted_at IS NULL`, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (r Repository) UpdateExistingUserById(ctx context.Context, id int, newUsername string, newPassword string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET username = $1, pass_hash = $2 WHERE id = $3`, newUsername, newPassword, id)

	return err
}

func (r Repository) UpdateUserRoleById(ctx context.Context, id int, newRole string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET role = $1 WHERE id = $2`, newRole, id)

	return err
}

func (r Repository) DeleteExistingUserById(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}

	if n.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r Repository) CheckIsUserExistsById(ctx context.Context, id int) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT count(*) FROM users WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
