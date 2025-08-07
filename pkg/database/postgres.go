package database

import (
	"database/sql"
	"fmt"

	"github.com/CAATHARSIS/movies-library/internal/config"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	conStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	return db, nil
}
