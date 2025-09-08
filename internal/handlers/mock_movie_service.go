package handlers

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/CAATHARSIS/movies-library/internal/models"
	"github.com/CAATHARSIS/movies-library/internal/service"
)

type MockMovieService struct {
	movies  map[int]*models.Movie
	nextID  int
	mu      sync.RWMutex
	ErrorOn map[string]bool
}

func NewMockMovieService() service.MovieService {
	return &MockMovieService{
		movies:  make(map[int]*models.Movie),
		nextID:  1,
		ErrorOn: make(map[string]bool),
	}
}

func (m *MockMovieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOn["CreateMovie"] {
		return errors.New("mock create movie error")
	}

	movie.ID = m.nextID
	m.nextID++
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	m.movies[movie.ID] = movie

	return nil
}

func (m *MockMovieService) GetMovie(ctx context.Context, id int) (*models.Movie, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOn["GetMovie"] {
		return nil, errors.New("mock get movie error")
	}

	movie, exists := m.movies[id]
	if !exists {
		return nil, errors.New("movie not found")
	}

	return movie, nil
}

func (m *MockMovieService) UpdateMovie(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOn["UpdateMovie"] {
		return nil, errors.New("mock update movie error")
	}

	old_movie, exists := m.movies[movie.ID]
	if !exists {
		return nil, errors.New("movie not found")
	}

	if movie.Title != "" {
		old_movie.Title = movie.Title
	}

	if movie.Director != "" {
		old_movie.Director = movie.Director
	}

	if movie.ReleaseDate.IsZero() {
		old_movie.ReleaseDate = movie.ReleaseDate
	}

	if movie.Genre != "" {
		old_movie.Genre = movie.Genre
	}

	if movie.Description != "" {
		old_movie.Description = movie.Description
	}

	old_movie.UpdatedAt = time.Now()

	return old_movie, nil
}

func (m *MockMovieService) DeleteMovie(ctx context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ErrorOn["DeleteMovie"] {
		return errors.New("mock delete movie error")
	}

	delete(m.movies, id)
	return nil
}

func (m *MockMovieService) ListMovies(ctx context.Context) ([]*models.Movie, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.ErrorOn["ListMovies"] {
		return nil, errors.New("mock list movies error")
	}

	var movies []*models.Movie
	for _, movie := range m.movies {
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MockMovieService) AddTestMovies(movies ...*models.Movie) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, movie := range movies {
		movie.ID = m.nextID
		m.nextID++
		movie.CreatedAt = time.Now()
		movie.UpdatedAt = time.Now()
		m.movies[movie.ID] = movie
	}
}

func (m *MockMovieService) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.movies = make(map[int]*models.Movie)
	m.nextID = 1
	m.ErrorOn = make(map[string]bool)
}

func (m *MockMovieService) SetErrorMode(methodName string, enable bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ErrorOn[methodName] = enable
}

func (m *MockMovieService) GetMovieCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.movies)
}

func (m *MockMovieService) GetMovieByID(id int) *models.Movie {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.movies[id]
}
