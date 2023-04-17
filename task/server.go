package task

import (
	"log"
	"net/http"
)

type Server struct {
	Port    string
	Handler http.Handler
}

func (s Server) Run() error {
	srv := &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Handler,
	}

	log.Println("starting http task server on port " + srv.Addr)

	return srv.ListenAndServe()
}
