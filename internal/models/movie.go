// Package models describes models which are using in application logic
package models

import "time"

// Movie struct describes relation in db
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	Genre       string    `json:"genre"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
