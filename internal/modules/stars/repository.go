package stars

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

func (r *Repository) GetStars(ctx context.Context) ([]*Star, error) {
	query, args, err := squirrel.Select("id", "first_name", "middle_name", "last_name", "avatar_url", "birth_date", "birth_place", "death_date", "bio", "created_at", "deleted_at").
		From("stars").
		Where(squirrel.Eq{"deleted_at": nil}).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	defer rows.Close()
	return r.scanStars(rows)
}

func (r *Repository) GetRelationsByMovieID(ctx context.Context, movieID int) ([]*MovieStarsRelation, error) {
	query, args, err := squirrel.Select("movie_id, star_id, hero_name, role, details, order_no").
		From("movie_stars").
		Where("movie_id = ?", movieID).
		OrderBy("order_no").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	rows, err := dbx.FromContext(ctx, r.db).Query(ctx, query, args...)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	var relations []*MovieStarsRelation
	for rows.Next() {
		var relation MovieStarsRelation
		if err = rows.Scan(&relation.MovieID, &relation.StarID, &relation.HeroName, &relation.Role, &relation.Details, &relation.OrderNo); err != nil {
			return nil, apperrors.Internal(err)
		}
		relations = append(relations, &relation)
	}

	return relations, nil
}

func (r *Repository) GetStarsPaginated(ctx context.Context, offset int, limit int) ([]*Star, int, error) {
	selectQuery := dbx.StatementBuilder.Select("id, first_name, middle_name, last_name, avatar_url, birth_date, birth_place, death_date, bio, created_at, deleted_at").
		From("stars").
		Where(squirrel.Eq{"deleted_at": nil}).
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	countQuery := dbx.StatementBuilder.Select("COUNT(*)").
		From("stars").
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

	stars, _ := r.scanStars(rows)

	if err = rows.Err(); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	var total int
	if err = br.QueryRow().Scan(&total); err != nil {
		return nil, 0, apperrors.Internal(err)
	}

	return stars, total, nil
}

func (r *Repository) GetStarByID(ctx context.Context, starID int) (*Star, error) {
	query, args, err := squirrel.Select("id", "first_name", "middle_name", "last_name", "avatar_url", "birth_date", "birth_place", "death_date", "bio", "imdb_url", "created_at").
		From("stars").
		Where(squirrel.Eq{"id": starID}).
		Where(squirrel.Eq{"deleted_at": nil}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	star := NewStar()
	err = r.db.QueryRow(ctx, query, args...).
		Scan(
			&star.ID,
			&star.FirstName,
			&star.MiddleName,
			&star.LastName,
			&star.AvatarURL,
			&star.BirthDate,
			&star.BirthPlace,
			&star.DeathDate,
			&star.Bio,
			&star.IMDbURL,
			&star.CreatedAt,
		)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("star", "id", starID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return star.Normalize(), nil
}

func (r *Repository) GetStarsForMovie(ctx context.Context, movieID int) ([]*Star, error) {
	query, args, err := squirrel.Select("id, first_name, middle_name, last_name, avatar_url, birth_date, birth_place, death_date, bio, created_at, deleted_at").
		From("stars").
		InnerJoin("movie_stars ON star_id = id").
		Where(squirrel.Eq{"movie_id": movieID}).
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

	stars, err := r.scanStars(rows)
	if err != nil {
		return nil, err
	} else if len(stars) == 0 {
		return nil, apperrors.NotFound("movie", "id", movieID)
	}

	return stars, nil
}

func (r *Repository) GetStarsByMovieID(ctx context.Context, movieID int) ([]*MovieCredit, error) {
	query, args, err := squirrel.Select("id, first_name, middle_name, last_name, avatar_url, birth_date, birth_place, death_date, imdb_url, bio, created_at, hero_name, role, details").
		From("stars").
		InnerJoin("movie_stars ON star_id = id").
		Where(squirrel.Eq{"movie_id": movieID}).
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

	var credits []*MovieCredit

	for rows.Next() {
		credit := &MovieCredit{
			Star: Star{},
		}

		err = rows.Scan(
			&credit.Star.ID,
			&credit.Star.FirstName,
			&credit.Star.MiddleName,
			&credit.Star.LastName,
			&credit.Star.AvatarURL,
			&credit.Star.BirthDate,
			&credit.Star.BirthPlace,
			&credit.Star.DeathDate,
			&credit.Star.IMDbURL,
			&credit.Star.Bio,
			&credit.Star.CreatedAt,
			&credit.HeroName,
			&credit.Role,
			&credit.Details,
		)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		credits = append(credits, credit)
	}

	return credits, nil
}

func (r *Repository) CreateStar(ctx context.Context, req *CreateStarRequest) (*Star, error) {
	query, args, err := squirrel.Insert("stars").
		Columns("first_name", "middle_name", "last_name", "avatar_url", "birth_date", "birth_place", "death_date", "bio", "imdb_url").
		Values(req.FirstName, req.MiddleName, req.LastName, req.AvatarURL, req.BirthDate, req.BirthPlace, req.DeathDate, req.Bio, req.IMDbURL).
		Suffix("RETURNING id, created_at").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	star := req.ToStar()
	err = r.db.QueryRow(ctx, query, args...).Scan(&star.ID, &star.CreatedAt)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return star.Normalize(), nil
}

func (r *Repository) UpdateStar(ctx context.Context, starID int, req *UpdateStarRequest) (*Star, error) {
	builder := squirrel.Update("stars").
		Where(squirrel.Eq{"id": starID}).
		Suffix("RETURNING id, first_name, middle_name, last_name, avatar_url, birth_date, birth_place, death_date, imdb_url, bio, created_at").
		PlaceholderFormat(squirrel.Dollar)

	hasSet := false
	if req.FirstName != nil {
		builder = builder.Set("first_name", *req.FirstName)
		hasSet = true
	}
	if req.MiddleName != nil {
		builder = builder.Set("middle_name", *req.MiddleName)
		hasSet = true
	}
	if req.LastName != nil {
		builder = builder.Set("last_name", *req.LastName)
		hasSet = true
	}
	if req.AvatarURL != nil {
		builder = builder.Set("avatar_url", *req.AvatarURL)
	}
	if req.BirthDate != nil {
		builder = builder.Set("birth_date", *req.BirthDate)
		hasSet = true
	}
	if req.BirthPlace != nil {
		builder = builder.Set("birth_place", *req.BirthPlace)
		hasSet = true
	}
	if req.DeathDate != nil {
		builder = builder.Set("death_date", *req.DeathDate)
		hasSet = true
	}
	if req.IMDbURL != nil {
		builder = builder.Set("imdb_url", *req.IMDbURL)
	}
	if req.Bio != nil {
		builder = builder.Set("bio", *req.Bio)
		hasSet = true
	}

	if !hasSet {
		builder = builder.Set("id", starID)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	star := NewStar()
	err = r.db.QueryRow(ctx, query, args...).Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.AvatarURL, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.IMDbURL, &star.Bio, &star.CreatedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("star", "id", starID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return star.Normalize(), nil
}

func (r *Repository) DeleteStarByID(ctx context.Context, starID int) error {
	query, args, err := squirrel.Update("stars").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": starID}).
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
		return apperrors.NotFound("star", "id", starID)
	}

	return nil
}

func (r *Repository) scanStars(rows pgx.Rows) ([]*Star, error) {
	var stars []*Star
	for rows.Next() {
		star := NewStar()
		err := rows.Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.AvatarURL, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt, &star.DeletedAt)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		stars = append(stars, star.Normalize())
	}
	return stars, nil
}
