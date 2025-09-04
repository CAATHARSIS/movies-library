package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/CAATHARSIS/movies-library/internal/models"
	"github.com/CAATHARSIS/movies-library/internal/service"
	"github.com/gorilla/mux"
)

type MovieHandler struct {
	service service.MovieService
	log     *slog.Logger
}

func NewMovieHandler(service service.MovieService, log *slog.Logger) *MovieHandler {
	return &MovieHandler{service: service, log: log}
}

func (h *MovieHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/movies", h.CreateMovie).Methods("POST")
	router.HandleFunc("/movies/{id}", h.GetMovie).Methods("GET")
	router.HandleFunc("/movies/{id}", h.UpdateMovie).Methods("PUT")
	router.HandleFunc("/movies/{id}", h.DeleteMovie).Methods("DELETE")
	router.HandleFunc("/movies", h.ListMovies).Methods("GET")
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Error("Failed to decode movie body", "error", err)
		return
	}

	if err := h.service.CreateMovie(r.Context(), &movie); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to create movie", "error", err)
		return
	}

	h.log.Info("Movie created succesully", "ID", movie.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid movie id", http.StatusBadRequest)
		h.log.Error("Invalid movie id", "ID", id)
		return
	}

	movie, err := h.service.GetMovie(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to get movie", "error", err)
		return
	}

	h.log.Info("Movie got", "ID", movie.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid movie id", http.StatusBadRequest)
		h.log.Error("Invalid movie ID", "error", err)
		return
	}

	var movie models.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Error("Failed to decode json body", "error", err)
		return
	}
	movie.ID = id

	updatedMovie, err := h.service.UpdateMovie(r.Context(), &movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to update movie", "ID", movie.ID)
		return
	}

	h.log.Info("Movie updated succesfully", "ID", movie.ID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedMovie)
}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid movie id", http.StatusBadRequest)
		h.log.Error("Invalid movie id", "error", err)
		return
	}

	if err := h.service.DeleteMovie(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to delete movie", "error", err)
		return
	}

	h.log.Info("Movie was deleted succesfully", "ID", err)
	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) ListMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.ListMovies(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to list movies", "error", err)
		return
	}

	h.log.Info("Movies listed succesfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
