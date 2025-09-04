package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CAATHARSIS/movies-library/internal/logger"
	"github.com/CAATHARSIS/movies-library/internal/models"
	"github.com/gorilla/mux"
)

func TestMovieHandler_CreateMovie_Succes(t *testing.T) {
	mockService := NewMockMovieService()
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	movie := &models.Movie{
		Title:       "Test Movie",
		Director:    "Test Director",
		ReleaseDate: time.Now(),
		Genre:       "Test Genre",
		Description: "Test Description",
	}

	body, _ := json.Marshal(movie)
	req := httptest.NewRequest("POST", "/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateMovie(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response models.Movie
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.ID == 0 {
		t.Error("Expected movie to have ID")
	}

	if response.Title != movie.Title {
		t.Errorf("Expected title %s, got %s", movie.Title, response.Title)
	}
}

func TestMovieHandler_CreateMovie_InvalidJSON(t *testing.T) {
	mockService := NewMockMovieService()
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	body := []byte(`{"invalid json": "something"`)
	req := httptest.NewRequest("POST", "/movies", bytes.NewReader(body))
	w := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/json")

	handler.CreateMovie(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestMovieHandler_CreateMovie_ServiceError(t *testing.T) {
	mockService := NewMockMovieService().(*MockMovieService)
	mockService.SetErrorMode("CreateMovie", true)
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	movie := &models.Movie{Title: "Test Movie"}

	body, _ := json.Marshal(movie)
	req := httptest.NewRequest("POST", "/movies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateMovie(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestMovieHandler_GetMovie_Succes(t *testing.T) {
	mockService := NewMockMovieService().(*MockMovieService)
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	testMovie := &models.Movie{
		Title: "Exsisting Movie",
		Director: "Test Director",
		ReleaseDate: time.Now(),
		Genre: "Test Genre",
	}

	mockService.AddTestMovies(testMovie)

	req := httptest.NewRequest("GET", "/movies/1", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/movies/{id}", handler.GetMovie).Methods("GET")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.Movie
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.ID != 1 {
		t.Errorf("Expected movie ID 1, got %d", response.ID)
	}

	if response.Title != "Exsisting Movie" {
		t.Errorf("Expected title 'Exsisting Movie', got '%s'", response.Title)
	}
}

func TestMovieHandler_GetMovie_NotFound(t *testing.T) {
	mockService := NewMockMovieService().(*MockMovieService)
	mockService.SetErrorMode("GetMovie", true)
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/movies/1", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/movies/{id}", handler.GetMovie).Methods("GET")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestMovieHandler_GetMovie_InvalidID(t *testing.T) {
	mockService := NewMockMovieService()
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/movies/invalid", nil)
	w := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/movies/{id}", handler.GetMovie).Methods("GET")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestMovieHandler_UpdateMovie_Partial(t *testing.T) {
	mockService := NewMockMovieService().(*MockMovieService)
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	initialMovie := &models.Movie{
		Title: "Original Title",
		Director: "Original Director",
		ReleaseDate: time.Now(),
		Genre: "Original Genre",
		Description: "Original Description",
	}

	mockService.AddTestMovies(initialMovie)

	updateData := map[string]interface{} {
		"Title": "Only new title",
	}

	body, _ := json.Marshal(updateData)

	req := httptest.NewRequest("PUT", "/movies/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/movies/{id}", handler.UpdateMovie).Methods("PUT")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	response := models.Movie{}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if response.Title != "Only new title" {
		t.Errorf("Expected title to be 'Only new title', got %s", response.Title)
	}

	if response.Director != "Original Director" {
		t.Error("Expected direcor would not be changed")
	}

	if response.Genre != "Original Genre" {
		t.Error("Expected genre would not be changed")
	}

	if response.Description != "Original Description" {
		t.Error("Expected description would not be changed")
	}
}

func TestMovieHandler_DeleteMovie_Succes(t *testing.T) {
	mockService := NewMockMovieService().(*MockMovieService)
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	initialMovie := &models.Movie{Title: "first movie"}

	mockService.AddTestMovies(initialMovie)

	req := httptest.NewRequest("DELETE", "/movies/1", nil)
	w := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/movies/{id}", handler.DeleteMovie).Methods("DELETE")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	if mockService.GetMovieCount() != 0 {
		t.Errorf("After deleting there must be one move, have %d", mockService.GetMovieCount())
	}
}

func TestMovieHandler_ListMovies_Succes(t *testing.T) {
	mockService := NewMockMovieService().(*MockMovieService)
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	initialMovies := []*models.Movie{
		&models.Movie{Title: "Title 1"},
		&models.Movie{Title: "Title 2"},
	}

	mockService.AddTestMovies(initialMovies...)

	req := httptest.NewRequest("GET", "/movies", nil)
	w := httptest.NewRecorder()
	handler.ListMovies(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response []*models.Movie

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if len(response) != 2 {
		t.Errorf("Expected movies count to be 2, got %d", len(response))
	}
}

func TestMovieHandler_ListMovies_Empty(t *testing.T) {
	mockService := NewMockMovieService()
	logger := logger.NewLogger("local")
	handler := NewMovieHandler(mockService, logger)

	req := httptest.NewRequest("GET", "/movies", nil)
	w := httptest.NewRecorder()
	handler.ListMovies(w, req)

	var response []*models.Movie

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	if len(response) != 0 {
		t.Errorf("Movie count must be zero, got %d", len(response))
	}
}