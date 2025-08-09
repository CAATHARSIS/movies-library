package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CAATHARSIS/movies-library/internal/config"
	"github.com/CAATHARSIS/movies-library/internal/handlers"
	"github.com/CAATHARSIS/movies-library/internal/repository/movie"
	"github.com/CAATHARSIS/movies-library/internal/service"
	"github.com/CAATHARSIS/movies-library/pkg/database"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	movieRepo := movie.NewMoviePostgresRepo(db)

	movieService := service.NewMovieService(movieRepo)

	movieHandler := handlers.NewMovieHandler(*movieService)

	router := mux.NewRouter()
	movieHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Server started on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<- done
	log.Println("Server is shutting down...")

	ctx, cancel :=  context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exitted properly")
}
