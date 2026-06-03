package app

import (
	"errors"
	"log"
	"log/slog"
	"phone-tokens/internal/service/billing"
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
	User        users.Service
	Token       tokens.Service
	SMS         *sms.SmsService
	Cert        *certificates.CertificateService
	Billing     *billing.BillingService
	UserProfile *users.UserProfileService
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
		slog.Error("failed to migrate database", "error", err)
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	userProfileRepo := repository.NewUserProfileRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	smsRepo := repository.NewSmsRepository(db)
	certificateRepo := repository.NewCertificateRepository(db)
	usageRepo := repository.NewUsageRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	pkgRepo := repository.NewPackageRepository(db)
	agentPkgRepo := repository.NewAgentPackageRepository(db)

	userSvc := users.NewService(userRepo, users.Config{
		JWTSecret:       cfg.JWTSecret,
		JWTExpiresInSec: cfg.JWTExpiresInSec,
	})

	tokenSvc := tokens.NewService(tokenRepo)
	userProfileSvc := users.NewUserProfileService(userRepo, userProfileRepo, tokenRepo)

	certSvc, err := certificates.NewCertificateService(certificateRepo)
	if err != nil {
		return nil, err
	}

	smsAdapter := sms_aero.NewAeroService(cfg.APIEmail, cfg.APIKey) // интерфейсную развязку сюда потом

	billingService := billing.NewBillingService(userRepo, usageRepo, transactionRepo,
		&pkgRepo, agentPkgRepo,
		cfg.BillingConfig.StripeKey, cfg.BillingConfig.WebhookSecret, cfg.FrontendURL)

	smsSvc := sms.NewSmsService(*certSvc, userProfileSvc, billingService, smsAdapter, tokenSvc, smsRepo)

	services := Services{
		User:        userSvc,
		Token:       tokenSvc,
		SMS:         smsSvc,
		Cert:        certSvc,
		Billing:     billingService,
		UserProfile: userProfileSvc,
	}
	return &services, nil
}
