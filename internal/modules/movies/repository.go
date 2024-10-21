package movies

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/DavidMovas/Movies-Reviews/internal/dbx"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"

	"github.com/DavidMovas/Movies-Reviews/contracts"
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

func (r *Repository) GetMovies(ctx context.Context, offset int, limit int, sort, order string) ([]*contracts.Movie, int, error) {
	b := &pgx.Batch{}
	b.Queue(`SELECT id, title, release_date, created_at, deleted_at FROM movies WHERE deleted_at IS NULL ORDER BY $1 $2 LIMIT $3 OFFSET $4`, sort, order, limit, offset)
	b.Queue(`SELECT COUNT(*) FROM movies WHERE deleted_at IS NULL`)

	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var movies []*contracts.Movie
	for rows.Next() {
		var movie contracts.Movie
		if err = rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseDate, &movie.CreatedAt, &movie.DeletedAt); err != nil {
			return nil, 0, apperrors.Internal(err)
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return movies, total, nil
}

func (r *Repository) GetMovieByID(ctx context.Context, movieID int) (*contracts.MovieDetails, error) {
	var movie contracts.MovieDetails

	err := r.db.QueryRow(ctx, `SELECT id, title, description, release_date, created_at, version FROM movies WHERE id = $1 AND deleted_at IS NULL`, movieID).
		Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.CreatedAt, &movie.Version)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("movie", "id", movieID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return &movie, nil
}

func (r *Repository) CreateMovie(ctx context.Context, req *contracts.CreateMovieRequest) (*contracts.MovieDetails, error) {
	var movie contracts.MovieDetails
	err := r.db.QueryRow(ctx, `INSERT INTO movies (title, description, release_date) VALUES ($1, $2, $3) RETURNING id, title, description, release_date, created_at, version`, req.Title, req.Description, req.ReleaseDate).
		Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.CreatedAt, &movie.Version)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return &movie, nil
}

func (r *Repository) UpdateMovieByID(ctx context.Context, movieID int, req *contracts.UpdateMovieRequest) (*contracts.MovieDetails, error) {
	fields := make(map[string]interface{})

	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.ReleaseDate != nil {
		fields["release_date"] = *req.ReleaseDate
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}

	var setClauses []string
	var values []interface{}
	index := 1

	for column, value := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, index))
		values = append(values, value)
		index++
	}

	query := fmt.Sprintf(`UPDATE movies SET %s, version = version + 1  WHERE id = $%d AND deleted_at IS NULL AND version = $%d  RETURNING id, title, description, release_date, created_at, deleted_at, version`, strings.Join(setClauses, ", "), index, req.Version)
	values = append(values, movieID)

	var movie contracts.MovieDetails
	err := r.db.QueryRow(ctx, query, values...).
		Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.CreatedAt, &movie.DeletedAt, &movie.Version)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("movie", "id", movieID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return &movie, nil
}

func (r *Repository) DeleteMovieByID(ctx context.Context, movieID int) error {
	n, err := r.db.Exec(ctx, `UPDATE movies SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, movieID)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("movie", "id", movieID)
	}

	return nil
}
