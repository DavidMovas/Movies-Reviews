package movies

import (
	"context"

	"github.com/Masterminds/squirrel"

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
	starsRepo  *stars.Repository
}

func NewRepository(db *pgxpool.Pool, genresRepo *genres.Repository, starsRepo *stars.Repository) *Repository {
	return &Repository{
		db:         db,
		genresRepo: genresRepo,
		starsRepo:  starsRepo,
	}
}

func (r *Repository) GetMovies(ctx context.Context, offset int, limit int, sort, order string, searchTerm *string) ([]*Movie, int, error) {
	selectQuery := dbx.StatementBuilder.Select("id, title, poster_url, release_date, avg_rating, created_at, deleted_at").
		From("movies").
		Where(squirrel.Eq{"deleted_at": nil}).
		OrderBy(sort + " " + order).
		Limit(uint64(limit)).
		Offset(uint64(offset))

	countQuery := dbx.StatementBuilder.Select("COUNT(*)").
		From("movies").
		Where(squirrel.Eq{"deleted_at": nil})

	if searchTerm != nil {
		selectQuery = selectQuery.Where("search_vector @@ to_tsquery('english', ?)", *searchTerm).
			OrderByClause("ts_rank_cd(search_vector, to_tsquery('english', ?)) DESC", *searchTerm)

		countQuery = countQuery.Where("search_vector @@ to_tsquery('english', ?)", *searchTerm)
	}

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

	var movies []*Movie
	for rows.Next() {
		var movie Movie
		if err = rows.Scan(&movie.ID, &movie.Title, &movie.PosterURL, &movie.ReleaseDate, &movie.AvgRating, &movie.CreatedAt, &movie.DeletedAt); err != nil {
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
	query, args, err := squirrel.Select("id", "title", "poster_url", "description", "imdb_rating", "imdb_url", "metascore", "metascore_url", "release_date", "avg_rating", "created_at", "version").
		From("movies").
		Where(squirrel.Eq{"id": movieID}).
		Where(squirrel.Eq{"deleted_at": nil}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	var movie MovieDetails
	err = r.db.QueryRow(ctx, query, args...).
		Scan(&movie.ID, &movie.Title, &movie.PosterURL, &movie.Description, &movie.IMDbRating, &movie.IMDbURL, &movie.Metascore, &movie.MetascoreURL, &movie.ReleaseDate, &movie.AvgRating, &movie.CreatedAt, &movie.Version)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("movie", "id", movieID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return &movie, nil
}

func (r *Repository) CreateMovie(ctx context.Context, movie *MovieDetails) error {
	query, args, err := squirrel.Insert("movies").
		Columns("title", "poster_url", "imdb_rating", "imdb_url", "metascore", "metascore_url", "description", "release_date").
		Values(movie.Title, movie.PosterURL, movie.IMDbRating, movie.IMDbURL, movie.Metascore, movie.MetascoreURL, movie.Description, movie.ReleaseDate).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return apperrors.Internal(err)
	}

	// Start transaction
	err = dbx.InTransaction(ctx, r.db, func(ctx context.Context, _ pgx.Tx) error {
		err = r.db.QueryRow(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt)
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

		// Insert stars
		nextCast := slices.MapIndex(movie.Cast, func(i int, credit *MovieCredit) *stars.MovieStarsRelation {
			if credit.HeroName == nil {
				credit.HeroName = ptr("")
			}
			return &stars.MovieStarsRelation{
				MovieID:  movie.ID,
				StarID:   credit.Star.ID,
				HeroName: *credit.HeroName,
				Role:     credit.Role,
				Details:  credit.Details,
				OrderNo:  i,
			}
		})
		return r.updateStars(ctx, nil, nextCast)
	})
	switch {
	case dbx.NotValidEnumType(err):
		return apperrors.BadRequestHidden(err, "role type unknown")
	case err != nil:
		return apperrors.Internal(err)
	}

	return nil
}

func (r *Repository) UpdateMovieByID(ctx context.Context, movieID int, req *UpdateMovieRequest) (*MovieDetails, error) {
	var movie MovieDetails
	builder := squirrel.Update("movies").
		Set("version", squirrel.Expr("version + 1")).
		Where(squirrel.Eq{"id": movieID}).
		Where(squirrel.Eq{"deleted_at": nil}).
		Where(squirrel.Eq{"version": req.Version}).
		Suffix("RETURNING id, title, poster_url, imdb_rating, imdb_url, metascore, metascore_url, description, release_date, created_at, deleted_at, version").
		PlaceholderFormat(squirrel.Dollar)

	if req.Title != nil {
		builder = builder.Set("title", *req.Title)
	}
	if req.ReleaseDate != nil {
		builder = builder.Set("release_date", *req.ReleaseDate)
	}
	if req.Description != nil {
		builder = builder.Set("description", *req.Description)
	}
	if req.PosterURL != nil {
		builder = builder.Set("poster_url", *req.PosterURL)
	}
	if req.IMDbRating != nil {
		builder = builder.Set("imdb_rating", *req.IMDbRating)
	}
	if req.IMDbURL != nil {
		builder = builder.Set("imdb_url", *req.IMDbURL)
	}
	if req.Metascore != nil {
		builder = builder.Set("metascore", *req.Metascore)
	}
	if req.MetascoreURL != nil {
		builder = builder.Set("metascore_url", *req.MetascoreURL)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	err = dbx.InTransaction(ctx, r.db, func(ctx context.Context, _ pgx.Tx) error {
		err = r.db.QueryRow(ctx, query, args...).
			Scan(
				&movie.ID,
				&movie.Title,
				&movie.PosterURL,
				&movie.IMDbRating,
				&movie.IMDbURL,
				&movie.Metascore,
				&movie.MetascoreURL,
				&movie.Description,
				&movie.ReleaseDate,
				&movie.CreatedAt,
				&movie.DeletedAt,
				&movie.Version,
			)

		switch {
		case dbx.NotValidEnumType(err):
			return apperrors.BadRequestHidden(err, "role type unknown")
		case dbx.IsNoRows(err):
			if _, err = r.GetMovieByID(ctx, movieID); err != nil {
				return apperrors.NotFound("movie", "id", movieID)
			} else {
				return apperrors.VersionMismatch("movie", "id", movieID, req.Version)
			}
		case err != nil:
			return apperrors.Internal(err)
		}

		if err = r.genresUpdateRequest(ctx, req.GenreIDs, movieID); err != nil {
			return err
		}

		if err = r.starsUpdateRequest(ctx, req.Cast, movieID); err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, apperrors.EnsureInternal(err)
	}

	return &movie, nil
}

func (r *Repository) DeleteMovieByID(ctx context.Context, movieID int) error {
	query, args, err := squirrel.Update("movies").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": movieID}).
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
		return apperrors.NotFound("movie", "id", movieID)
	}

	return nil
}

func (r *Repository) Lock(ctx context.Context, tx pgx.Tx, movieID int) error {
	n, err := tx.Exec(ctx, `SELECT 1 FROM movies WHERE deleted_at IS NULL AND id = $1 FOR UPDATE`, movieID)
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
		_, err := q.Exec(ctx, `INSERT INTO movie_stars (movie_id, star_id, hero_name, role, details, order_no) VALUES ($1, $2, $3, $4, $5, $6)`,
			mgo.MovieID, mgo.StarID, mgo.HeroName, mgo.Role, mgo.Details, mgo.OrderNo)
		return err
	}

	removeFunc := func(mgo *stars.MovieStarsRelation) error {
		_, err := q.Exec(ctx, `DELETE FROM movie_stars WHERE movie_id = $1 AND star_id = $2 AND hero_name = $3 AND role = $4 AND details = $5`, mgo.MovieID, mgo.StarID, mgo.HeroName, mgo.Role, mgo.Details)
		return err
	}

	return dbx.AdjustRelation(current, next, addFunc, removeFunc)
}

func (r *Repository) genresUpdateRequest(ctx context.Context, ids []*int, movieID int) error {
	if ids != nil {
		var err error
		var genresIDs []int
		for _, genre := range ids {
			genresIDs = append(genresIDs, *genre)
		}

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

	return nil
}

func (r *Repository) starsUpdateRequest(ctx context.Context, cast []*MovieCreditInfo, movieID int) error {
	if cast != nil {
		var err error
		var starsIDs []int
		for _, star := range cast {
			starsIDs = append(starsIDs, star.StarID)
		}

		var currenCast []*stars.MovieStarsRelation
		currenCast, err = r.starsRepo.GetRelationsByMovieID(ctx, movieID)
		if err != nil {
			return err
		}

		nextCast := slices.MapIndex(starsIDs, func(i, starID int) *stars.MovieStarsRelation {
			if cast[i].HeroName == nil {
				cast[i].HeroName = ptr("N/A")
			}
			return &stars.MovieStarsRelation{
				MovieID:  movieID,
				StarID:   starID,
				HeroName: *cast[i].HeroName,
				Role:     cast[i].Role,
				Details:  cast[i].Details,
				OrderNo:  i,
			}
		})

		return r.updateStars(ctx, currenCast, nextCast)
	}
	return nil
}

func ptr(ptr string) *string {
	return &ptr
}
