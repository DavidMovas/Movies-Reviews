package users

import (
	"context"

	"github.com/DavidMovas/Movies-Reviews/contracts"
	"github.com/DavidMovas/Movies-Reviews/internal/dbx"
	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
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

func (r Repository) Create(ctx context.Context, user *contracts.UserWithPassword) (err error) {
	err = r.db.QueryRow(ctx,
		`INSERT INTO users (username, email, pass_hash, role) 
    		VALUES ($1, $2, $3, $4) RETURNING id, role, created_at`, user.Username, user.Email, user.PasswordHash, user.Role).
		Scan(&user.ID, &user.Role, &user.CreatedAt)

	switch {
	case dbx.IsUniqueViolation(err, "email"):
		return apperrors.AlreadyExists("user", "email", user.Email)
	case dbx.IsUniqueViolation(err, "username"):
		return apperrors.AlreadyExists("user", "username", user.Username)
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (r Repository) GetExistingUserByEmail(ctx context.Context, email string) (*contracts.UserWithPassword, error) {
	user := contracts.NewUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role, created_at, deleted_at FROM users WHERE email = $1 AND deleted_at IS NULL`, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return user, nil
}

func (r Repository) GetExistingUserById(ctx context.Context, id int) (*contracts.UserWithPassword, error) {
	user := contracts.NewUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role, created_at, deleted_at FROM users WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", id)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return user, nil
}

func (r Repository) GetExistingUserByUsername(ctx context.Context, username string) (*contracts.UserWithPassword, error) {
	user := contracts.NewUserWithPassword()
	err := r.db.QueryRow(ctx, `SELECT id, username, email, pass_hash, role, created_at FROM users WHERE username = $1 AND deleted_at IS NULL`, username).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return user, nil
}

func (r Repository) UpdateExistingUserById(ctx context.Context, id int, newUsername string, newPassword string) error {
	n, err := r.db.Exec(ctx, `UPDATE users SET username = $1, pass_hash = $2 WHERE id = $3`, newUsername, newPassword, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}

	return nil
}

func (r Repository) UpdateUserRoleById(ctx context.Context, id int, newRole string) error {
	n, err := r.db.Exec(ctx, `UPDATE users SET role = $1 WHERE id = $2`, newRole, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}

	return nil
}

func (r Repository) DeleteExistingUserById(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, `UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}

	return nil
}

func (r Repository) CheckIsUserExistsById(ctx context.Context, id int) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT count(*) FROM users WHERE id = $1 AND deleted_at IS NULL`, id).
		Scan(&count)
	if err != nil {
		return false, apperrors.Internal(err)
	}

	if count == 0 {
		return false, apperrors.NotFound("user", "id", id)
	}

	return count > 0, nil
}
