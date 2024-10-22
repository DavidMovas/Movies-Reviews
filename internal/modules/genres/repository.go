package genres

import (
	"context"

	"github.com/jackc/pgx/v5"

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

func (r *Repository) GetGenres(ctx context.Context) ([]*Genre, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM genres`)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	defer rows.Close()

	return scanGenres(rows)
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, movieID int) ([]*MovieGenreRelation, error) {
	rows, err := dbx.FromContext(ctx, r.db).
		Query(ctx, `SELECT movie_id, genre_id, order_no FROM movie_genres WHERE movie_id = $1 ORDER BY order_no`, movieID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	var relations []*MovieGenreRelation
	for rows.Next() {
		var relation MovieGenreRelation
		if err = rows.Scan(&relation.MovieID, &relation.GenreID, &relation.OrderNo); err != nil {
			return nil, apperrors.Internal(err)
		}
		relations = append(relations, &relation)
	}
	return relations, nil
}

func (r *Repository) GetGenreByID(ctx context.Context, id int) (*Genre, error) {
	var genre Genre
	err := r.db.QueryRow(ctx, `SELECT id, name FROM genres WHERE id = $1`, id).
		Scan(&genre.ID, &genre.Name)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("genre", "id", id)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return &genre, nil
}

func (r *Repository) GetGenresByMovieID(ctx context.Context, movieID int) ([]*Genre, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name FROM genres INNER JOIN movie_genres ON genre_id = id WHERE movie_id = $1 ORDER BY order_no`, movieID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	return scanGenres(rows)
}

func (r *Repository) CreateGenre(ctx context.Context, raq *CreateGenreRequest) (*Genre, error) {
	var genre Genre
	err := r.db.QueryRow(ctx, `INSERT INTO genres(name) VALUES($1) ON CONFLICT (name) DO NOTHING  RETURNING id, name`, raq.Name).
		Scan(&genre.ID, &genre.Name)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.AlreadyExists("genre", "name", raq.Name)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return &genre, nil
}

func (r *Repository) UpdateGenreByID(ctx context.Context, id int, raq *UpdateGenreRequest) error {
	n, err := r.db.Exec(ctx, `UPDATE genres SET name = $1 WHERE id = $2 
        AND NOT EXISTS (SELECT 1 FROM genres WHERE name = $1 AND id <> $2)`, raq.Name, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		var exists bool
		err = r.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM genres WHERE name = $1)`, raq.Name).Scan(&exists)
		if err != nil {
			return apperrors.Internal(err)
		}

		if exists {
			return apperrors.AlreadyExists("genre", "name", raq.Name)
		}

		return apperrors.NotFound("genre", "id", id)
	}

	return nil
}

func (r *Repository) DeleteGenreByID(ctx context.Context, genreID int) error {
	n, err := r.db.Exec(ctx, `DELETE FROM genres WHERE id = $1`, genreID)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", genreID)
	}

	return nil
}

func scanGenres(rows pgx.Rows) ([]*Genre, error) {
	var genres []*Genre
	for rows.Next() {
		var genre Genre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, apperrors.Internal(err)
		}
		genres = append(genres, &genre)
	}

	return genres, nil
}
