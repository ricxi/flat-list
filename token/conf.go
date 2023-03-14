package token

import (
	"errors"
	"os"
)

type Conf struct {
	DatabaseURL string
}

func GetConf() (*Conf, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, errors.New("db connection env cannot be empty")
	}

	return &Conf{DatabaseURL: connStr}, nil
}
