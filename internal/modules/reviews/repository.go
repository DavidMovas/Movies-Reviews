package reviews

import (
	"context"
	"fmt"
	"time"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"
	"github.com/jackc/pgx/v5"

	"github.com/Masterminds/squirrel"

	"github.com/DavidMovas/Movies-Reviews/internal/dbx"

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

func (r *Repository) GetReviewsByMovieID(ctx context.Context, movieID int, offset int, limit int, sort, order string) ([]*Review, int, error) {
	selectQuery := dbx.StatementBuilder.Select("id, movie_id, user_id, rating, title, content, created_at, updated_at, deleted_at").
		From("reviews").
		Where("movie_id = ?", movieID).
		Where(squirrel.Eq{"deleted_at": nil}).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		OrderBy(sort + " " + order)

	countQuery := dbx.StatementBuilder.Select("COUNT(*)").
		From("reviews").
		Where("movie_id = ?", movieID).
		Where(squirrel.Eq{"deleted_at": nil})

	b := &pgx.Batch{}
	if err := dbx.QueryBatchSelect(b, selectQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	if err := dbx.QueryBatchSelect(b, countQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	br := r.db.SendBatch(ctx, b)
	defer func() {
		_ = br.Close()
	}()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows[*Review](rows, pgx.RowToAddrOfStructByPos[Review])
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return reviews, total, nil
}

func (r *Repository) GetReviewsByUserID(ctx context.Context, userID int, offset int, limit int, sort, order string) ([]*Review, int, error) {
	selectQuery := dbx.StatementBuilder.Select("id, movie_id, user_id, rating, title, content, created_at, updated_at, deleted_at").
		From("reviews").
		Where("user_id = ?", userID).
		Where(squirrel.Eq{"deleted_at": nil}).
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		OrderBy(sort + " " + order)

	countQuery := dbx.StatementBuilder.Select("COUNT(*)").
		From("reviews").
		Where("user_id = ?", userID).
		Where(squirrel.Eq{"deleted_at": nil})

	b := &pgx.Batch{}
	if err := dbx.QueryBatchSelect(b, selectQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	if err := dbx.QueryBatchSelect(b, countQuery); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	br := r.db.SendBatch(ctx, b)
	defer func() {
		_ = br.Close()
	}()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	reviews, err := pgx.CollectRows[*Review](rows, pgx.RowToAddrOfStructByPos[Review])
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return reviews, total, nil
}

func (r *Repository) GetReviewByID(ctx context.Context, reviewID int) (*Review, error) {
	builder := dbx.StatementBuilder.Select("id, movie_id, user_id, rating, title, content, created_at, updated_at, deleted_at").
		From("reviews").
		Where("id = ?", reviewID).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("review", "id", reviewID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	return pgx.RowToAddrOfStructByPos[Review](rows)
}

func (r *Repository) CreateReview(ctx context.Context, req *CreateReviewRequest) (*Review, error) {
	builder := dbx.StatementBuilder.Insert("reviews").
		Columns("movie_id", "user_id", "rating", "title", "content").
		Values(req.MovieID, req.UserID, req.Rating, req.Title, req.Content).
		Suffix("RETURNING id, rating, title, content, created_at, updated_at, deleted_at")

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var review Review
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&review.ID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt, &review.DeletedAt)

	switch {
	case dbx.IsUniqueViolation(err, "reviews_movie_id_user_id_key"):
		return nil, apperrors.AlreadyExists("review", "movie_id user_id", fmt.Sprintf("%d - %d", req.MovieID, req.UserID))
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	review.MovieID = req.MovieID
	review.UserID = req.UserID

	return &review, nil
}

func (r *Repository) UpdateReview(ctx context.Context, reviewID int, req *UpdateReviewRequest) (*Review, error) {
	if req.Rating == nil && req.Title == nil && req.Content == nil {
		return nil, apperrors.BadRequest(fmt.Errorf("no fields to update"))
	}

	builder := dbx.StatementBuilder.Update("reviews").
		Set("updated_at", time.Now()).
		Where("id = ?", reviewID).
		Where(squirrel.Eq{"deleted_at": nil}).
		Suffix("RETURNING id, movie_id, user_id, rating, title, content, created_at, updated_at, deleted_at")

	if req.Rating != nil {
		builder = builder.Set("rating", *req.Rating)
	}
	if req.Title != nil {
		builder = builder.Set("title", *req.Title)
	}
	if req.Content != nil {
		builder = builder.Set("content", *req.Content)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var review Review
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&review.ID, &review.MovieID, &review.UserID, &review.Rating, &review.Title, &review.Content, &review.CreatedAt, &review.UpdatedAt, &review.DeletedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("review", "id", reviewID)
	case err != nil:
		return nil, apperrors.Internal(err)
	}

	return &review, nil
}

func (r *Repository) DeleteReview(ctx context.Context, reviewID int) error {
	builder := dbx.StatementBuilder.Update("reviews").
		Set("deleted_at", time.Now()).
		Where("id = ?", reviewID).
		Where(squirrel.Eq{"deleted_at": nil})

	query, args, err := builder.ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	n, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("review", "id", reviewID)
	}

	return nil
}
