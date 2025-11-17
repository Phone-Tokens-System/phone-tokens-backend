package app

import (
	"errors"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	httpadapter "users/internal/users/adapter/in/http"
	"users/internal/users/adapter/out/repository"
	"users/internal/users/model"
	"users/internal/users/service/users"
)

func NewHTTPServer(cfg Config) (*http.Server, error) {
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.HTTPPort == "" {
		return nil, errors.New("HTTP_PORT is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	if cfg.JWTExpiresInSec <= 0 {
		return nil, errors.New("JWT_EXPIRES_IN_SEC must be greater than zero")
	}

	log.Printf("initializing PostgreSQL repository and running migrations")
	repo, err := NewPostgresRepository(cfg.DatabaseURL)
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

func NewPostgresRepository(databaseURL string) (users.Repository, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations on startup to ensure schema is up to date.
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return repository.NewPostgresRepository(db), nil
}
