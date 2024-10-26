package users

import (
	"context"

	"github.com/Masterminds/squirrel"

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

func (r Repository) Create(ctx context.Context, user *UserWithPassword) (err error) {
	query, args, err := squirrel.Insert("users").
		Columns("username", "email", "pass_hash", "role").
		Values(user.Username, user.Email, user.PasswordHash, user.Role).
		Suffix("RETURNING id, role, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return apperrors.Internal(err)
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Role, &user.CreatedAt)

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

func (r Repository) GetExistingUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	query, args, err := squirrel.Select("id, username, email, pass_hash, role, created_at, deleted_at").
		From("users").
		Where("email = ? AND deleted_at IS NULL", email).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := NewUserWithPassword()
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return user, nil
}

func (r Repository) GetExistingUserByID(ctx context.Context, id int) (*UserWithPassword, error) {
	query, args, err := squirrel.Select("id, username, email, pass_hash, role, created_at, deleted_at").
		From("users").
		Where("id = ? AND deleted_at IS NULL", id).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := NewUserWithPassword()
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", id)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return user, nil
}

func (r Repository) GetExistingUserByUsername(ctx context.Context, username string) (*UserWithPassword, error) {
	query, args, err := squirrel.Select("id, username, email, pass_hash, role, created_at, deleted_at").
		From("users").
		Where("username = ? AND deleted_at IS NULL", username).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := NewUserWithPassword()
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return user, nil
}

func (r Repository) UpdateExistingUserByID(ctx context.Context, id int, newUsername string, newPassword string) error {
	query, args, err := squirrel.Update("users").
		Set("username", newUsername).
		Set("pass_hash", newPassword).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return apperrors.Internal(err)
	}

	n, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}

	return nil
}

func (r Repository) UpdateUserRoleByID(ctx context.Context, id int, newRole string) error {
	query, args, err := squirrel.Update("users").
		Set("role", newRole).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return apperrors.Internal(err)
	}

	n, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}

	return nil
}

func (r Repository) DeleteExistingUserByID(ctx context.Context, id int) error {
	query, args, err := squirrel.Update("users").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return apperrors.Internal(err)
	}

	n, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("user", "id", id)
	}

	return nil
}
