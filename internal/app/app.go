package app

import (
	"errors"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"phone_token_system/internal/adapter/out/repository"
	"phone_token_system/internal/service/tokens"
	"phone_token_system/internal/service/users"
)

func BuildService(cfg Config) (users.Service, tokens.Service, error) {
	if cfg.DatabaseURL == "" {
		return nil, nil, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		return nil, nil, errors.New("JWT_SECRET is required")
	}
	if cfg.JWTExpiresInSec <= 0 {
		return nil, nil, errors.New("JWT_EXPIRES_IN_SEC must be greater than zero")
	}

	log.Printf("initializing database connection")
	db, errInit := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if errInit != nil {
		return nil, nil, errInit
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	log.Printf("applying database migrations with goose")
	if err := MigrateDB(sqlDB, "database"); err != nil {
		return nil, nil, err
	}

	repo := repository.NewStorage(db)

	userSvc := users.NewService(repo, users.Config{
		JWTSecret:       cfg.JWTSecret,
		JWTExpiresInSec: cfg.JWTExpiresInSec,
	})

	tokenSvc := tokens.NewService(repo)

	return userSvc, tokenSvc, nil
}
