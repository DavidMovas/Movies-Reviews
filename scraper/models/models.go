package models

import "time"

type Movie struct {
	Title       string
	Description string
	Genres      []*Genre
	ReleaseDate time.Time
}

type Genre struct {
	Name string
}
