package movie

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CAATHARSIS/movies-library/internal/models"
)

type MoviePostgresRepo struct {
	db *sql.DB
}

func NewMoviePostgresRepo(db *sql.DB) *MoviePostgresRepo {
	return &MoviePostgresRepo{db}
}

func (r *MoviePostgresRepo) Create(ctx context.Context, movie *models.Movie) error {
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

func (r *MoviePostgresRepo) GetByID(ctx context.Context, id int) (*models.Movie, error) {
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

func (r *MoviePostgresRepo) Update(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
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

	var updatedMovie models.Movie
	err := r.db.QueryRowContext(
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

func (r *MoviePostgresRepo) Delete(ctx context.Context, id int) error {
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

func (r *MoviePostgresRepo) List(ctx context.Context) ([]*models.Movie, error) {
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
