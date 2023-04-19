package user

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type server struct {
	srv *http.Server
}

func NewServer(handler http.Handler, port string) *server {
	return &server{
		&http.Server{
			Addr:         ":" + port,
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
	}
}

func (s *server) Run() {
	shutdownErr := make(chan error, 1)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

		receiver := <-sigs
		log.Println("received signal to start shutdown:", receiver)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shutdownErr <- s.srv.Shutdown(ctx)
	}()

	log.Printf("starting user service on port %s...\n", s.srv.Addr)

	err := s.srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		// TODO: check for other errors
		log.Fatalln("expecting http.ErrServerClosed to indicate successful graceful shutdown", err)
	}

	if err := <-shutdownErr; err != nil {
		log.Fatalln("problem during shutdown", err)
	}
	log.Printf("service on port %s has been shut down", s.srv.Addr)
}
