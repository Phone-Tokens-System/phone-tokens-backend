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
	NovofonAPIKey   string
	NovofonSecret   string
	NovofonBaseURL  string
	NovofonTimeout  int64
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file %w", err)
	}

	novofonAPIKey := os.Getenv("NOVOFON_API_KEY")
	if novofonAPIKey == "" {
		novofonAPIKey = os.Getenv("ZADARMA_API_KEY")
	}

	novofonSecret := os.Getenv("NOVOFON_API_SECRET")
	if novofonSecret == "" {
		novofonSecret = os.Getenv("ZADARMA_API_SECRET")
	}

	novofonBaseURL := os.Getenv("NOVOFON_BASE_URL")
	if novofonBaseURL == "" {
		novofonBaseURL = os.Getenv("ZADARMA_BASE_URL")
	}

	novofonTimeout := getenvInt64("NOVOFON_TIMEOUT_SEC")
	if novofonTimeout == 0 {
		novofonTimeout = getenvInt64("ZADARMA_TIMEOUT_SEC")
	}

	return Config{
		HTTPPort:        os.Getenv("HTTP_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiresInSec: getenvInt64("JWT_EXPIRES_IN_SEC"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		APIKey:          os.Getenv("API_KEY"),
		APIEmail:        os.Getenv("EMAIL"),
		NovofonAPIKey:   novofonAPIKey,
		NovofonSecret:   novofonSecret,
		NovofonBaseURL:  novofonBaseURL,
		NovofonTimeout:  novofonTimeout,
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
