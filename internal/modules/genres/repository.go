package genres

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

func (r *Repository) GetGenres(ctx context.Context) ([]*contracts.Genre, error) {
	var genres []*contracts.Genre
	rows, err := r.db.Query(ctx, `SELECT id, name FROM genres`)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	defer rows.Close()

	for rows.Next() {
		genre := contracts.Genre{}
		err = rows.Scan(&genre.Name, &genre.Name)
		if err != nil {
			return nil, apperrors.InternalWithoutStackTrace(err)
		}
		genres = append(genres, &genre)
	}
	return genres, nil
}

func (r *Repository) GetGenreById(ctx context.Context, id int) (*contracts.Genre, error) {
	var genre *contracts.Genre
	err := r.db.QueryRow(ctx, `SELECT id, name FROM genres WHERE id = $1`, id).
		Scan(&genre.ID, &genre.Name)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("genre", "id", id)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return genre, nil
}

func (r *Repository) CreateGenre(ctx context.Context, raq *contracts.CreateGenreRequest) (*contracts.Genre, error) {
	var genre *contracts.Genre
	err := r.db.QueryRow(ctx, `INSERT INTO genres(name) VALUES($1) RETURNING id, name`, raq.Name).
		Scan(&genre.ID, &genre.Name)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("genre", "name", raq.Name)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}
	return genre, nil
}

func (r *Repository) UpdateGenreById(ctx context.Context, id int, raq *contracts.UpdateGenreRequest) error {
	n, err := r.db.Exec(ctx, `UPDATE genres SET name = $1 WHERE id = $2`, raq.Name, id)

	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", id)
	}

	return nil
}

func (r *Repository) DeleteGenreById(ctx context.Context, id int) error {
	n, err := r.db.Exec(ctx, `DELETE FROM genres WHERE id = $1`, id)
	if err != nil {
		return apperrors.Internal(err)
	}

	if n.RowsAffected() == 0 {
		return apperrors.NotFound("genre", "id", id)
	}

	return nil
}
