package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var ErrInvalidEnv = errors.New("invalid environment variable")

type MissingEnvErr struct {
	missingEnvs []string
}

func (m *MissingEnvErr) Error() string {
	return fmt.Sprintf("missing environment variables: %v", m.missingEnvs)
}

func (m *MissingEnvErr) add(missing string) {
	m.missingEnvs = append(m.missingEnvs, missing)
}

func (m *MissingEnvErr) hasErrors() bool {
	return len(m.missingEnvs) != 0
}

type envs struct {
	port         string
	mongoURI     string
	mongoDBName  string
	mongoTimeout int
}

func LoadEnvs() (*envs, error) {
	errs := MissingEnvErr{}

	port := os.Getenv("PORT")
	if port == "" {
		errs.add("PORT")
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		errs.add("MONGODB_URI")
	}

	mongoDBName := os.Getenv("MONGODB_NAME")
	if mongoDBName == "" {
		errs.add("MONGODB_NAME")
	}

	mongoTimeoutStr := os.Getenv("MONGODB_TIMEOUT")
	if mongoTimeoutStr == "" {
		errs.add("MONGODB_TIMEOUT")
	}
	mongoTimeout, err := strconv.Atoi(mongoTimeoutStr)
	if err != nil {
		return nil, fmt.Errorf("%w: MONGODB_TIMEOUT", ErrInvalidEnv)
	}

	if errs.hasErrors() {
		return nil, &errs
	}

	return &envs{
		port:         port,
		mongoURI:     mongoURI,
		mongoDBName:  mongoDBName,
		mongoTimeout: mongoTimeout,
	}, nil
}
