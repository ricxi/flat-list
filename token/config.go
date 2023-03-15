package token

import (
	"errors"
	"os"
)

type config struct {
	DatabaseURL string
}

func LoadConfig() (*config, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, errors.New("db connection env cannot be empty")
	}

	return &config{DatabaseURL: connStr}, nil
}
