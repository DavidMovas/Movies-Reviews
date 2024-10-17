package contracts

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CreateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=32"`
}

type UpdateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=32"`
}

func NewGenre() *Genre {
	return &Genre{}
}
