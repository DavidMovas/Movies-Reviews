package models

import "time"

type Movie struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	PosterURL   string    `json:"poster_url"`
	IMDbRating  float64   `json:"imdb_rating"`
	Description string    `json:"description"`
	Metascore   int       `json:"metascore"`
	Storyline   string    `json:"storyline"`
	Runtime     string    `json:"runtime"`
	Genres      []string  `json:"genres"`
	ReleaseDate time.Time `json:"release_date"`

	Link string `json:"_link"`
}

type Genre struct {
	Name string
}

type Cast struct {
	MovieID string    `json:"movie_id"`
	Cast    []*Credit `json:"cast"`

	Link string `json:"_link"`
}

type Credit struct {
	Role     string `json:"role"`
	Details  string `json:"details"`
	StarID   string `json:"star_id"`
	StarName string `json:"star_name"`
	HeroName string `json:"hero_name"`
	StarLink string `json:"_star_link"`
}

type Star struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	FirstName  string     `json:"first_name"`
	MiddleName *string    `json:"middle_name"`
	LastName   string     `json:"last_name"`
	AvatarURL  string     `json:"avatar_url"`
	IMDbURL    string     `json:"imdb_url"`
	BirthDate  time.Time  `json:"birth_date"`
	DeathDate  *time.Time `json:"death_date"`

	Link string `json:"_link"`
}

type Bio struct {
	ID         string `json:"id"`
	Bio        string `json:"bio"`
	BirthPlace string `json:"birth_place"`

	Link string `json:"_link"`
}
