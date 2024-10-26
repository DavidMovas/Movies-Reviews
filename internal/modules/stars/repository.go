package stars

import (
	"context"
	"fmt"
	"strings"

	"github.com/DavidMovas/Movies-Reviews/internal/modules/movies"

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
	var stars []*Star
	rows, err := r.db.Query(ctx, `SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at FROM stars WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	defer rows.Close()

	for rows.Next() {
		star := NewStar()
		err = rows.Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt)
		if err != nil {
			return nil, apperrors.InternalWithoutStackTrace(err)
		}
		stars = append(stars, star.Normalize())
	}
	return stars, nil
}

func (r *Repository) GetStarsPaginated(ctx context.Context, offset int, limit int) ([]*Star, int, error) {
	b := &pgx.Batch{}
	b.Queue(`SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at, deleted_at FROM stars WHERE deleted_at IS NULL ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	b.Queue(`SELECT COUNT(*) FROM stars WHERE deleted_at IS NULL`)
	br := r.db.SendBatch(ctx, b)
	defer br.Close()

	rows, err := br.Query()
	if err != nil {
		return nil, 0, apperrors.Internal(err)
	}
	defer rows.Close()

	var stars []*Star
	stars, err = r.scanStars(rows)

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
	star := NewStar()
	err := r.db.QueryRow(ctx, `SELECT id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at FROM stars WHERE id = $1 AND deleted_at IS NULL`, starID).
		Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("star", "id", starID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return star.Normalize(), nil
}

func (r *Repository) GetStarsByMovieID(ctx context.Context, movieID int) ([]*movies.MovieCredit, error) {
	rows, err := r.db.Query(ctx, `SELECT s.id, s.first_name, s.middle_name, s.last_name, s.birth_date, s.birth_place, s.death_date, s.bio, s.created_at, ms.role, ms.details
						FROM stars s INNER JOIN movie_stars ms ON star_id = id WHERE movie_id = $1 ORDER BY order_no`, movieID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	defer rows.Close()

	var credits []*movies.MovieCredit

	for rows.Next() {
		credit := &movies.MovieCredit{
			Star: Star{},
		}

		err = rows.Scan(&credit.Star.ID, &credit.Star.FirstName, &credit.Star.MiddleName, &credit.Star.LastName, &credit.Star.BirthDate, &credit.Star.BirthPlace, &credit.Star.DeathDate, &credit.Star.Bio, &credit.Star.CreatedAt, &credit.Role, &credit.Details)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		credits = append(credits, credit)
	}

	return credits, nil
}

func (r *Repository) CreateStar(ctx context.Context, req *CreateStarRequest) (*Star, error) {
	star := req.ToStar()
	err := r.db.QueryRow(ctx, `INSERT INTO stars (first_name, middle_name, last_name, birth_date, birth_place, death_date, bio) 
									VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`,
		req.FirstName, req.MiddleName, req.LastName, req.BirthDate, req.BirthPlace, req.DeathDate, req.Bio).
		Scan(&star.ID, &star.CreatedAt)
	if err != nil {
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return star.Normalize(), nil
}

func (r *Repository) UpdateStar(ctx context.Context, starID int, req *UpdateStarRequest) (*Star, error) {
	fields := make(map[string]interface{})

	if req.FirstName != nil {
		fields["first_name"] = *req.FirstName
	}
	if req.MiddleName != nil {
		fields["middle_name"] = *req.MiddleName
	}
	if req.LastName != nil {
		fields["last_name"] = *req.LastName
	}
	if req.BirthDate != nil {
		fields["birth_date"] = *req.BirthDate
	}
	if req.BirthPlace != nil {
		fields["birth_place"] = *req.BirthPlace
	}
	if req.DeathDate != nil {
		fields["death_date"] = *req.DeathDate
	}
	if req.Bio != nil {
		fields["bio"] = *req.Bio
	}

	var setClauses []string
	var values []interface{}
	index := 1

	for column, value := range fields {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, index))
		values = append(values, value)
		index++
	}

	query := fmt.Sprintf(`UPDATE stars SET %s WHERE id = $%d RETURNING id, first_name, middle_name, last_name, birth_date, birth_place, death_date, bio, created_at`, strings.Join(setClauses, ", "), index)
	values = append(values, starID)

	star := NewStar()
	err := r.db.QueryRow(ctx, query, values...).Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt)

	switch {
	case dbx.IsNoRows(err):
		return nil, apperrors.NotFound("star", "id", starID)
	case err != nil:
		return nil, apperrors.InternalWithoutStackTrace(err)
	}

	return star.Normalize(), nil
}

func (r *Repository) DeleteStarByID(ctx context.Context, starID int) error {
	n, err := r.db.Exec(ctx, `UPDATE stars SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`, starID)
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
		err := rows.Scan(&star.ID, &star.FirstName, &star.MiddleName, &star.LastName, &star.BirthDate, &star.BirthPlace, &star.DeathDate, &star.Bio, &star.CreatedAt, &star.DeletedAt)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
		stars = append(stars, star.Normalize())
	}
	return stars, nil
}
