package app

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort        string
	JWTSecret       string
	JWTExpiresInSec int64
	DatabaseURL     string
}

func LoadConfig() Config {
	return Config{
		HTTPPort:        os.Getenv("HTTP_PORT"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		JWTExpiresInSec: getenvInt64("JWT_EXPIRES_IN_SEC"),
		DatabaseURL:     os.Getenv("DATABASE_URL"),
	}
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
