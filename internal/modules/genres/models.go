package genres

import "github.com/DavidMovas/Movies-Reviews/internal/dbx"

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var _ dbx.Keyer = MovieGenreRelation{}

type MovieGenreRelation struct {
	MovieID int
	GenreID int
	OrderNo int
}

func (m MovieGenreRelation) Key() any {
	type MovieGenreRelationKey struct {
		MovieID, GenreID int
	}

	return MovieGenreRelationKey{
		MovieID: m.MovieID,
		GenreID: m.GenreID,
	}
}

type CreateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=32"`
}

type UpdateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=32"`
}
