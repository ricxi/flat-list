package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
)

// var wg sync.WaitGroup
// wg.Add(2)
func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tokenDirGRPC := filepath.Join(wd, "../token/cmd/grpc")
	tokenDirHTTP := filepath.Join(wd, "../token/cmd/http")

	tokenEnvs := os.Environ()
	tokenEnvs = append(
		tokenEnvs,
		"DATABASE_URL=postgres://postgres:password@127.0.0.1:5433/tokens",
		"HTTP_PORT=5002",
		"GRPC_PORT=5003",
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

	go httpTokenSvc.run(&wg)
	go grpcTokenSvc.run(&wg)

	wg.Wait()
}
