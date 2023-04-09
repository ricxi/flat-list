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

type DBVars struct {
	PsqlDSN string `mapstructure:"PSQL_DSN"`
}

type config struct {
	ServiceVars `mapstructure:"services"`
	DBVars      `mapstructure:"databases"`
}

func main() {
	var c config
	filename := "config"
	if err := LoadTOMLConfig("", filename, &c); err != nil {
		log.Fatalln("cannot start services without configuration variables", err)
	}

	scripts := [][]string{
		{"./start_mongo.sh"},
		{"./start_postgres.sh", c.PsqlDSN},
	}
	// scriptErr := make(chan error)
	var wgs sync.WaitGroup
	wgs.Add(len(scripts))
	for _, script := range scripts {
		go func(script ...string) {
			defer wgs.Done()
			if err := runSH(script...); err != nil {
				// scriptErr <- err
				log.Println(err)
			}
		}(script...)
	}
	wgs.Wait()

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
