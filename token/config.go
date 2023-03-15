package token

import (
	"errors"
	"os"
)

type config struct {
	DatabaseURL string
	GrpcPort    string
	HttpPort    string
}

func LoadConfig() (*config, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, errors.New("db connection env cannot be empty")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return nil, errors.New("HTTP_PORT environment variable could not be found")
	}

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		return nil, errors.New("GRPC_PORT environment variable could not be found")
	}

	return &config{
		DatabaseURL: connStr,
		HttpPort:    httpPort,
		GrpcPort:    grpcPort,
	}, nil
}
