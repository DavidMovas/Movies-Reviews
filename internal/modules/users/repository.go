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
		Columns("username", "email", "pass_hash", "role, avatar_url").
		Values(user.Username, user.Email, user.PasswordHash, user.Role, user.AvatarURL).
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
	query, args, err := squirrel.Select("id, username, email, pass_hash, role, avatar_url, bio, created_at, deleted_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		Where(squirrel.Eq{"deleted_at": nil}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := NewUserWithPassword()
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "email", email)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return user, nil
}

func (r Repository) GetExistingUserByID(ctx context.Context, id int) (*User, error) {
	query, args, err := squirrel.Select("id, username, email, role, avatar_url, bio, created_at, deleted_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Eq{"deleted_at": nil}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var user User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "id", id)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return &user, nil
}

func (r Repository) GetExistingUserByUsername(ctx context.Context, username string) (*UserWithPassword, error) {
	query, args, err := squirrel.Select("id, username, email, role, avatar_url, bio, created_at, deleted_at").
		From("users").
		Where(squirrel.Eq{"username": username}).
		Where(squirrel.Eq{"deleted_at": nil}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	user := NewUserWithPassword()
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("user", "username", username)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return user, nil
}

func (r Repository) UpdateExistingUserByID(ctx context.Context, id int, req *UpdateUserRequest, newPassword string) (*User, error) {
	builder := squirrel.Update("users").
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING id, username, email, role, avatar_url, bio, created_at, deleted_at").
		PlaceholderFormat(squirrel.Dollar)

	hasSet := false
	if req.Username != nil {
		builder = builder.Set("username", *req.Username)
		hasSet = true
	}
	if newPassword != "" {
		builder = builder.Set("pass_hash", newPassword)
		hasSet = true
	}
	if req.Bio != nil {
		builder = builder.Set("bio", *req.Bio)
		hasSet = true
	}
	if req.AvatarURL != nil {
		builder = builder.Set("avatar_url", *req.AvatarURL)
		hasSet = true
	}

	if !hasSet {
		builder = builder.Set("id", id)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var user User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.Bio, &user.CreatedAt, &user.DeletedAt)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	return &user, nil
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
