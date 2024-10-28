package genres

import "github.com/DavidMovas/Movies-Reviews/internal/dbx"

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetGenreRequest struct {
	GenreID int `json:"-" param:"genreId" validate:"nonzero"`
}

type CreateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=32"`
}

type UpdateGenreRequest struct {
	GenreID int    `json:"-" param:"genreId" validate:"nonzero"`
	Name    string `json:"name" validate:"min=3,max=32"`
}

type DeleteGenreRequest struct {
	GenreID int `json:"-" param:"genreId" validate:"nonzero"`
}

// MovieGenreRelation is a relation between movie and genre
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
