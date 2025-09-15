// Package service provides communication between repository and handler
package service

import (
	"context"

	"github.com/CAATHARSIS/movies-library/internal/models"
	"github.com/CAATHARSIS/movies-library/internal/repository/movie"
)

// MovieService interface describes structs that are used for creating handlers
type MovieService interface {
	CreateMovie(context.Context, *models.Movie) error
	GetMovie(context.Context, int) (*models.Movie, error)
	UpdateMovie(context.Context, *models.Movie) (*models.Movie, error)
	DeleteMovie(context.Context, int) error
	ListMovies(context.Context) ([]*models.Movie, error)
}

type movieService struct {
	repo movie.Repository
}

// NewMovieService creates new instance of MovieService interface
func NewMovieService(r movie.Repository) MovieService {
	return &movieService{repo: r}
}

func (s *movieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	return s.repo.Create(ctx, movie)
}

func (s *movieService) GetMovie(ctx context.Context, id int) (*models.Movie, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *movieService) UpdateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	return s.repo.Update(ctx, movie)
}

func (s *movieService) DeleteMovie(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *movieService) ListMovies(ctx context.Context) ([]*models.Movie, error) {
	return s.repo.List(ctx)
}
