package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort        string
	JWTSecret       string
	JWTExpiresInSec int64
	DatabaseURL     string
	APIKey          string
	APIEmail        string
	FrontendURL     string // e.g. "http://localhost:5173"
	BillingConfig   BillingConfig
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file %v", err)
	}
	billingConfig := BillingConfig{
		StripeKey:     os.Getenv("STRIPE_KEY"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		ServerURL:     os.Getenv("SERVER_URL"),
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	return Config{
		HTTPPort:        os.Getenv("HTTP_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiresInSec: getenvInt64("JWT_EXPIRES_IN_SEC"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		APIKey:          os.Getenv("API_KEY"),
		APIEmail:        os.Getenv("EMAIL"),
		FrontendURL:     frontendURL,
		BillingConfig:   billingConfig,
	}, nil
}

func getenvInt64(key string) int64 {
	value := os.Getenv(key)
	if value == "" {
		return 0
	}

	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0
	}

	return n
}

type BillingConfig struct {
	ServerURL     string // базовый URL, например http://localhost:8080
	StripeKey     string // секретный ключ Stripe
	WebhookSecret string // секрет для webhook
}
