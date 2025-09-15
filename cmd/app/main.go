package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CAATHARSIS/movies-library/internal/config"
	"github.com/CAATHARSIS/movies-library/internal/handlers"
	"github.com/CAATHARSIS/movies-library/internal/logger"
	"github.com/CAATHARSIS/movies-library/internal/middleware"
	"github.com/CAATHARSIS/movies-library/internal/repository/movie"
	"github.com/CAATHARSIS/movies-library/internal/service"
	"github.com/CAATHARSIS/movies-library/pkg/database"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	log := logger.NewLogger(cfg.Env)

	log.Info("starting movies-library", slog.String("env", cfg.Env))
	if cfg.Env != "prod" {
		log.Info("debug messages are enabled")
	}

	migrationDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.RunMigrations(migrationDB, log); err != nil {
		log.Error("Failed to run migrations", "error", err)
		if err := migrationDB.Close(); err != nil {
			log.Error("Failed to close migration db", "error", err)
		}
		os.Exit(1)
	}

	appDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Error("Failed to connect to database", "error", err)
		if err := appDB.Close(); err != nil {
			log.Error("Failed to close app db", "error", err)
		}
		os.Exit(1)
	}

	movieRepo := movie.NewMoviePostgresRepo(appDB)

	movieService := service.NewMovieService(movieRepo)

	movieHandler := handlers.NewMovieHandler(movieService, log)

	router := mux.NewRouter()
	router.Use(middleware.NewLoggingMiddleware(log))
	movieHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Server started", "port", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
			if err := migrationDB.Close(); err != nil {
				log.Error("Failed to close migration db", "error", err)
			}
			if err := appDB.Close(); err != nil {
				log.Error("Failed to close app db", "error", err)
			}
			os.Exit(1)
		}
	}()

	<-done
	log.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown faild", "error", err)
		if err := migrationDB.Close(); err != nil {
			log.Error("Failed to close migration db", "error", err)
		}
		if err := appDB.Close(); err != nil {
			log.Error("Failed to close app db", "error", err)
		}
		return
	}

	if err := migrationDB.Close(); err != nil {
		log.Error("Failed to close migration db", "error", err)
	}
	if err := appDB.Close(); err != nil {
		log.Error("Failed to close app db", "error", err)
	}

	log.Info("Server exitted properly")
}
