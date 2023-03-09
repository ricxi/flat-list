package main

import (
	"errors"
	"os"
	"strconv"
)

var ErrMissingEnvs = errors.New("missing environment variable")

type config struct {
	port         string
	mongoURI     string
	mongoDBName  string
	mongoTimeout int
}

func getEnvs() (config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return config{}, ErrMissingEnvs
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return config{}, ErrMissingEnvs
	}

	mongoDBName := os.Getenv("MONGODB_NAME")
	if mongoDBName == "" {
		return config{}, ErrMissingEnvs
	}

	mongoTimeoutStr := os.Getenv("MONGODB_TIMEOUT")
	if mongoTimeoutStr == "" {
		return config{}, ErrMissingEnvs
	}
	mongoTimeout, err := strconv.Atoi(mongoTimeoutStr)
	if err != nil {
		return config{}, nil
	}

	return config{
		port:         port,
		mongoURI:     mongoURI,
		mongoDBName:  mongoDBName,
		mongoTimeout: mongoTimeout,
	}, nil
}
