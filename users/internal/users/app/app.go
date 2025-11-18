package app

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"users/internal/users/adapter/out/repository"
	"users/internal/users/model"
	"users/internal/users/service/users"
)

func BuildService(cfg Config) (users.Service, error) {
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}
	if cfg.JWTExpiresInSec <= 0 {
		return nil, errors.New("JWT_EXPIRES_IN_SEC must be greater than zero")
	}

	log.Printf("initializing database connection")
	db, errInit := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if errInit != nil {
		return nil, errInit
	}

	log.Printf("applying database migrations")
	if errMigrate := db.AutoMigrate(&model.User{}); errMigrate != nil {
		return nil, errMigrate
	}

	repo := repository.NewStorage(db)

	return users.NewService(repo, users.Config{
		JWTSecret:       cfg.JWTSecret,
		JWTExpiresInSec: cfg.JWTExpiresInSec,
	}), nil
}
