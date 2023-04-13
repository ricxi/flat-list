package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
)

type Service struct {
	Name        string   `mapstructure:"name"`
	Dir         string   `mapstructure:"dir"`
	Envs        []string `mapstructure:"envs"`
	InitScripts []string `mapstructure:"init_scripts"`
}

// Services stores the configuration settings
// for each service that will be run.
type Services struct {
	Token  Service `mapstructure:"token"`
	Mailer Service `mapstructure:"mailer"`
	User   Service `mapstructure:"user"`
}

type config struct {
	Services `mapstructure:"services"`
}

func main() {
	var c config
	filename := "config"
	if err := LoadTOMLConfig("", filename, &c); err != nil {
		log.Fatalln("cannot start services without configuration variables", err)
	}

	mailerSvc, err := newServiceRunner(c.Mailer)
	if err != nil {
		log.Fatalln(err)
	}
	tokenSvc, err := newServiceRunner(c.Token)
	if err != nil {
		log.Fatalln(err)
	}
	userSvc, err := newServiceRunner(c.User)
	if err != nil {
		log.Fatalln(err)
	}

	services := []serviceRunner{
		mailerSvc,
		tokenSvc,
		userSvc,
	}
	// This is used to cancel all the services.
	// It's cancel function is called as soon as a signal is received on the signal channel (see below)
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, svc := range services {
		go func(svc serviceRunner) {
			defer wg.Done()
			if err := svc.initialize(); err != nil {
				return
			}

			// this should block if successful
			if err := svc.run(ctx); err != nil {
				if !errors.As(err, new(*exec.ExitError)) {
					log.Println("unexpected error ", err)
				}
				return
			}
		}(svc)
	}

	sChan := make(chan os.Signal, 1)
	signal.Notify(sChan, os.Interrupt, syscall.SIGTERM)
	log.Println("Exiting with signal: ", <-sChan)
	cancel()

	wg.Wait()
}
