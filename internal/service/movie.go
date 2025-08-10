package service

import (
	"context"

	"github.com/CAATHARSIS/movies-library/internal/models"
	"github.com/CAATHARSIS/movies-library/internal/repository/movie"
)

type MovieService struct {
	repo *movie.MoviePostgresRepo
}

func NewMovieService(r *movie.MoviePostgresRepo) *MovieService {
	return &MovieService{repo: r}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	return s.repo.Create(ctx, movie)
}

func (s *MovieService) GetMovie(ctx context.Context, id int) (*models.Movie, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	return s.repo.Update(ctx, movie)
}

func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *MovieService) ListMovies(ctx context.Context) ([]*models.Movie, error) {
	return s.repo.List(ctx)
}
