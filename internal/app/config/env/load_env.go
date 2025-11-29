package env

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

func LoadConfigEnv() (*CredsEnv, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file %w", err)
	}
	creds := CredsEnv{}
	err = env.Parse(&creds)
	if err != nil {
		return nil, fmt.Errorf("error parsing .env file %w", err)
	}

	return &creds, nil
}
