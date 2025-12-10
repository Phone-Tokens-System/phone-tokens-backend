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
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file %w", err)
	}

	return Config{
		HTTPPort:        os.Getenv("HTTP_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiresInSec: getenvInt64("JWT_EXPIRES_IN_SEC"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		APIKey:          os.Getenv("API_KEY"),
		APIEmail:        os.Getenv("EMAIL"),
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
