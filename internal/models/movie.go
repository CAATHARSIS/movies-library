package models

import "time"

type Movie struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Director string `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	Gengre string `json:"gengre"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}