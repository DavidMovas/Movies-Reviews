package movies

import (
	"context"
	"fmt"
	"strings"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/stars"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/genres"
	"github.com/DavidMovas/Movies-Reviews/internal/slices"

	"github.com/jackc/pgx/v5"

	"github.com/DavidMovas/Movies-Reviews/internal/dbx"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db         *pgxpool.Pool
	genresRepo *genres.Repository
}

func NewRepository(db *pgxpool.Pool, genresRepo *genres.Repository) *Repository {
	return &Repository{
		db:         db,
		genresRepo: genresRepo,
	}
}

func (r *Repository) GetMovies(ctx context.Context, offset int, limit int, sort, order string) ([]*Movie, int, error) {
	query := fmt.Sprintf(`SELECT id, title, release_date, created_at, deleted_at 
    FROM movies 
    WHERE deleted_at IS NULL 
    ORDER BY %s %s LIMIT $1 OFFSET $2`, sort, order)

	b := &pgx.Batch{}
	b.Queue(query, limit, offset)
	b.Queue(`SELECT COUNT(*) FROM movies WHERE deleted_at IS NULL`)
	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var movies []*Movie
	for rows.Next() {
		var movie Movie
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

func (r *Repository) GetMovieByID(ctx context.Context, movieID int) (*MovieDetails, error) {
	var movie MovieDetails

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

func (r *Repository) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, _ pgx.Tx) error {
		err := r.db.QueryRow(ctx, `INSERT INTO movies (title, description, release_date) VALUES ($1, $2, $3) RETURNING id, created_at`, movie.Title, movie.Description, movie.ReleaseDate).
			Scan(&movie.ID, &movie.CreatedAt)
		if err != nil {
			return err
		}

		// Insert genres
		nextGenres := slices.MapIndex(movie.Genres, func(i int, genre *genres.Genre) *genres.MovieGenreRelation {
			return &genres.MovieGenreRelation{
				MovieID: movie.ID,
				GenreID: genre.ID,
				OrderNo: i,
			}
		})
		if err = r.updateGenres(ctx, nil, nextGenres); err != nil {
			return err
		}

		nextCast := slices.MapIndex(movie.Cast, func(i int, credit *MovieCredit) *stars.MovieStarsRelation {
			return &stars.MovieStarsRelation{
				MovieID: movie.ID,
				StarID:  credit.Star.ID,
				Role:    credit.Role,
				Details: credit.Details,
				OrderNo: i,
			}
		})
		return r.updateStars(ctx, nil, nextCast)
	})
	if err != nil {
		return apperrors.Internal(err)
	}

	return nil
}

func (r *Repository) UpdateMovieByID(ctx context.Context, movieID int, req *UpdateMovieRequest) (*MovieDetails, error) {
	var movie MovieDetails

	var genresIDs []int
	for _, genre := range req.GenreIDs {
		genresIDs = append(genresIDs, *genre)
	}

	err := dbx.InTransaction(ctx, r.db, func(ctx context.Context, _ pgx.Tx) error {
		query, values := r.prepareQueryForUpdateRequest(movieID, req)
		err := r.db.QueryRow(ctx, query, values...).
			Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseDate, &movie.CreatedAt, &movie.DeletedAt, &movie.Version)

		switch {
		case dbx.IsNoRows(err):
			if _, err = r.GetMovieByID(ctx, movieID); err != nil {
				return apperrors.NotFound("movie", "id", movieID)
			} else {
				return apperrors.VersionMismatch("movie", "id", movieID, req.Version)
			}
		case err != nil:
			return apperrors.Internal(err)
		}

		if req.GenreIDs != nil {
			var currentGenres []*genres.MovieGenreRelation
			currentGenres, err = r.genresRepo.GetRelationsByMovieID(ctx, movieID)
			if err != nil {
				return err
			}

			nextGenres := slices.MapIndex(genresIDs, func(i, genreID int) *genres.MovieGenreRelation {
				return &genres.MovieGenreRelation{
					MovieID: movieID,
					GenreID: genreID,
					OrderNo: i,
				}
			})

			return r.updateGenres(ctx, currentGenres, nextGenres)
		}

		return err
	})
	if err != nil {
		return nil, apperrors.EnsureInternal(err)
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

func (r *Repository) updateGenres(ctx context.Context, current, next []*genres.MovieGenreRelation) error {
	q := dbx.FromContext(ctx, r.db)

	addFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx, `INSERT INTO movie_genres (movie_id, genre_id, order_no) VALUES ($1, $2, $3)`, mgo.MovieID, mgo.GenreID, mgo.OrderNo)
		return err
	}

	removeFunc := func(mgo *genres.MovieGenreRelation) error {
		_, err := q.Exec(ctx, `DELETE FROM movie_genres WHERE movie_id = $1 AND genre_id = $2`, mgo.MovieID, mgo.GenreID)
		return err
	}

	return dbx.AdjustRelation(current, next, addFunc, removeFunc)
}

func (r *Repository) updateStars(ctx context.Context, current, next []*stars.MovieStarsRelation) error {
	q := dbx.FromContext(ctx, r.db)

	addFunc := func(mgo *stars.MovieStarsRelation) error {
		_, err := q.Exec(ctx, `INSERT INTO movies_stars (movie_id, star_id, role, details, order_no) VALUES ($1, $2, $3, $4, $5)`,
			mgo.MovieID, mgo.StarID, mgo.Role, mgo.Details, mgo.OrderNo)
		return err
	}

	removeFunc := func(mgo *stars.MovieStarsRelation) error {
		_, err := q.Exec(ctx, `DELETE FROM movies_stars WHERE movie_id = $1 AND star_id = $2 AND role = $3 AND details = $4`, mgo.MovieID, mgo.StarID, mgo.Role, mgo.Details)
		return err
	}

	return dbx.AdjustRelation(current, next, addFunc, removeFunc)
}

func (r *Repository) prepareQueryForUpdateRequest(movieID int, req *UpdateMovieRequest) (string, []any) {
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

	query := fmt.Sprintf(`UPDATE movies SET %s, version = version + 1 WHERE id = $%d AND deleted_at IS NULL AND version = $%d  RETURNING id, title, description, release_date, created_at, deleted_at, version`, strings.Join(setClauses, ", "), index, index+1)
	values = append(values, movieID, req.Version)

	return query, values
}
