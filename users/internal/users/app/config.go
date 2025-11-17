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
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	expStr := os.Getenv("JWT_EXPIRES_IN_SEC")
	if expStr == "" {
		expStr = "3600"
	}

	exp, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		exp = 3600
	}

	return Config{
		HTTPPort:        port,
		JWTSecret:       secret,
		JWTExpiresInSec: exp,
		DatabaseURL:     os.Getenv("DATABASE_URL"),
	}
}
