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

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// I should put the setup for the token service into its
	// own function and use command-line flags to give the
	// option of setting up either one or both services
	tokenDirGRPC := filepath.Join(wd, "../token/cmd/grpc")
	tokenDirHTTP := filepath.Join(wd, "../token/cmd/http")

	tokenEnvs := append(
		os.Environ(),
		c.TokenVars...,
	)

	// scripts that should be run before
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

	errChan := runInitScripts("./start_postgres.sh", c.PsqlDSN)
	err = <-errChan
	if err != nil {
		log.Println()
		os.Exit(1)
	} else {
		var wg sync.WaitGroup
		wg.Add(3)

		go httpTokenSvc.run(&wg)
		go grpcTokenSvc.run(&wg)

		// run scripts to initialize the mongo user database
		go func() {
			defer wg.Done()

			if err := runSH("./start_mongo.sh"); err != nil {
				log.Println(err)
				return
			}
		}()
		wg.Wait()
	}
}
