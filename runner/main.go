package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

// ServiceVars stores a slice of
// configuration variables for
// each service that is run
type ServiceVars struct {
	TokenVars  []string `mapstructure:"token"`
	MailerVars []string `mapstructure:"mailer"`
	UserVars   []string `mapstructure:"user"`
}

type config struct {
	ServiceVars `mapstructure:"services"`
}

func main() {
	var c config
	filename := "config"
	if err := LoadTOMLConfig("", filename, &c); err != nil {
		log.Fatalln("cannot start services without configuration variables", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// This should be a flag
	tokenDirGRPC := filepath.Join(wd, "../token/cmd/grpc")
	tokenDirHTTP := filepath.Join(wd, "../token/cmd/http")

	tokenEnvs := append(
		os.Environ(),
		c.TokenVars...,
	)

	httpTokenSvc := goService{
		name:    "http token",
		workDir: tokenDirHTTP,
		envs:    tokenEnvs,
	}

	grpcTokenSvc := goService{
		name:    "grpc token",
		workDir: tokenDirGRPC,
		envs:    tokenEnvs,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// These errors aren't being caught
	go httpTokenSvc.run(&wg)
	go grpcTokenSvc.run(&wg)

	wg.Wait()
}
