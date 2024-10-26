package genres

import (
	"context"

	"github.com/Masterminds/squirrel"

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
	query, args, err := squirrel.Select("id, name").
		From("genres").
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	defer rows.Close()

	return scanGenres(rows)
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, movieID int) ([]*MovieGenreRelation, error) {
	query, args, err := squirrel.Select("movie_id, genre_id, order_no").
		From("movie_genres").
		Where("movie_id = $1", movieID).
		OrderBy("order_no").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	rows, err := dbx.FromContext(ctx, r.db).
		Query(ctx, query, args...)
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
	query, args, err := squirrel.Select("id, name").
		From("genres").
		Where("id = $1", id).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var genre Genre
	err = r.db.QueryRow(ctx, query, args...).
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
	query, args, err := squirrel.Select("id, name").
		From("genres").
		InnerJoin("movie_genres ON genre_id = id").
		Where("movie_id = $1", movieID).
		OrderBy("order_no").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	return scanGenres(rows)
}

func (r *Repository) CreateGenre(ctx context.Context, raq *CreateGenreRequest) (*Genre, error) {
	query, args, err := squirrel.Insert("genres(name)").
		Values(raq.Name).
		Suffix("ON CONFLICT (name) DO NOTHING RETURNING id, name").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var genre Genre
	err = r.db.QueryRow(ctx, query, args...).
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
	query, args, err := squirrel.Update("genres").
		Set("name = $1", raq.Name).
		Where("id = $2 AND NOT EXISTS (SELECT 1 FROM genres WHERE name = $1 AND id <> $2) ", id).
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
	query, args, err := squirrel.Delete("genres").
		Where("id = $1", genreID).
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
