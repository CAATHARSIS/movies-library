// Package movie provides communication application with db
package movie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CAATHARSIS/movies-library/internal/models"
)

// Repository interface describes functions which object must implements to communicate with db
type Repository interface {
	Create(context.Context, *models.Movie) error
	GetByID(context.Context, int) (*models.Movie, error)
	Update(context.Context, *models.Movie) (*models.Movie, error)
	Delete(context.Context, int) error
	List(context.Context) ([]*models.Movie, error)
}

type moviePostgresRepo struct {
	db *sql.DB
}

// NewMoviePostgresRepo creates new instance of moviePostgresRepo
func NewMoviePostgresRepo(db *sql.DB) Repository {
	return &moviePostgresRepo{db}
}

func (r *moviePostgresRepo) Create(ctx context.Context, movie *models.Movie) error {
	qurery := `
		INSERT INTO
			movies (
				title,
				director,
				release_date,
				genre,
				description,
				created_at,
				updated_at
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id
	`

	err := r.db.QueryRowContext(
		ctx,
		qurery,
		movie.Title,
		movie.Director,
		movie.ReleaseDate,
		movie.Genre,
		movie.Description,
		time.Now(),
		time.Now(),
	).Scan(&movie.ID)

	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

func (r *moviePostgresRepo) GetByID(ctx context.Context, id int) (*models.Movie, error) {
	query := `
		SELECT
			id,
			title,
			director,
			release_date,
			genre,
			description,
			created_at,
			updated_at
		FROM
			movies
		WHERE
			id = $1
	`

	var movie models.Movie

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Director,
		&movie.ReleaseDate,
		&movie.Genre,
		&movie.Description,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("movie not found")
		}
		return nil, fmt.Errorf("failed to get movie: %v", err)
	}

	return &movie, nil
}

func (r *moviePostgresRepo) Update(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	query := `
		UPDATE
			movies
		SET title = $1,
			director = $2,
			release_date = $3,
			genre = $4,
			description = $5,
			updated_at = $6
		WHERE
			id = $7
		RETURNING
			id,
			title,
			director,
			release_date,
			genre,
			description,
			created_at,
			updated_at
	`

	oldMovie, err := r.GetByID(ctx, movie.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid movie id: %v", err)
	}

	if movie.Title == "" {
		movie.Title = oldMovie.Title
	}

	if movie.Director == "" {
		movie.Director = oldMovie.Director
	}

	if movie.ReleaseDate.IsZero() {
		movie.ReleaseDate = oldMovie.ReleaseDate
	}

	if movie.Genre == "" {
		movie.Genre = oldMovie.Genre
	}

	if movie.Description == "" {
		movie.Description = oldMovie.Description
	}

	var updatedMovie models.Movie
	err = r.db.QueryRowContext(
		ctx,
		query,
		movie.Title,
		movie.Director,
		movie.ReleaseDate,
		movie.Genre,
		movie.Description,
		time.Now(),
		movie.ID,
	).Scan(
		&updatedMovie.ID,
		&updatedMovie.Title,
		&updatedMovie.Director,
		&updatedMovie.ReleaseDate,
		&updatedMovie.Genre,
		&updatedMovie.Description,
		&updatedMovie.CreatedAt,
		&updatedMovie.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update movie: %v", err)
	}

	return &updatedMovie, nil
}

func (r *moviePostgresRepo) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM movies
		WHERE
			id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete movie: %v", err)
	}

	return nil
}

func (r *moviePostgresRepo) List(ctx context.Context) ([]*models.Movie, error) {
	query := `
		SELECT
			ID,
			TITLE,
			DIRECTOR,
			RELEASE_DATE,
			GENRE,
			DESCRIPTION,
			CREATED_AT,
			UPDATED_AT
		FROM
			MOVIES
		ORDER BY
			UPDATED_AT DESC
	`

	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to list movies: %v", err)
	}

	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Director,
			&movie.ReleaseDate,
			&movie.Genre,
			&movie.Description,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %v", err)
		}

		movies = append(movies, &movie)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return movies, nil
}
