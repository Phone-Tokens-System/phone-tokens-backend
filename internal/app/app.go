package app

import (
	"errors"
	"log"
	"phone-tokens/internal/service/certificates"
	"phone-tokens/internal/service/sms"
	"phone-tokens/internal/service/sms/sms_aero"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"phone-tokens/internal/adapter/out/repository"
	"phone-tokens/internal/service/tokens"
	"phone-tokens/internal/service/users"
)

type Services struct {
	User  users.Service
	Token tokens.Service
	SMS   *sms.SmsService
	Cert  *certificates.CertificateService
}

func BuildService(cfg Config) (*Services, error) {
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

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	log.Printf("applying database migrations with goose")
	if err := MigrateDB(sqlDB, "database"); err != nil {
		return nil, err
	}

	repo := repository.NewStorage(db)

	userSvc := users.NewService(repo, users.Config{
		JWTSecret:       cfg.JWTSecret,
		JWTExpiresInSec: cfg.JWTExpiresInSec,
	})

	tokenSvc := tokens.NewService(repo)

	certSvc, err := certificates.NewCertificateService(repo)
	if err != nil {
		return nil, err
	}

	smsAdapter := sms_aero.NewAeroService(cfg.APIEmail, cfg.APIKey) // интерфейсную развязку сюда потом

	smsSvc := sms.NewSmsService(*certSvc, smsAdapter, tokenSvc, repo)

	services := Services{
		User:  userSvc,
		Token: tokenSvc,
		SMS:   smsSvc,
		Cert:  certSvc,
	}
	return &services, nil
}
