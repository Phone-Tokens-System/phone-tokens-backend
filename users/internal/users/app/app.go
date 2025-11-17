package app

import (
	"errors"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	httpadapter "users/internal/users/adapter/in/http"
	"users/internal/users/adapter/out/repository"
	"users/internal/users/service/users"
)

func NewHTTPServer(cfg Config) (*http.Server, error) {
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	log.Printf("initializing PostgreSQL repository")
	repo, err := newPostgresRepository(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	svc := users.NewService(repo, users.Config{
		JWTSecret:       cfg.JWTSecret,
		JWTExpiresInSec: cfg.JWTExpiresInSec,
	})

	handler := httpadapter.NewHandler(svc)

	mux := http.NewServeMux()
	httpadapter.RegisterRoutes(mux, handler, httpadapter.AuthConfig{
		JWTSecret: cfg.JWTSecret,
	})

	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: mux,
	}

	log.Printf("HTTP server initialized on port %s", cfg.HTTPPort)

	return server, nil
}

func newPostgresRepository(databaseURL string) (users.Repository, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return repository.NewPostgresRepository(db), nil
}
