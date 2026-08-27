package api

import (
	"log"
	"net/http"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	s.registerRoutes(mux)

	srv := &http.Server{
		Addr:    s.addr,
		Handler: mux,
	}

	log.Printf("Starting server on %s", s.addr)
	return srv.ListenAndServe()
}

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", s.GetUsersHandler)
	mux.HandleFunc("POST /users", s.CreateUserHandler)
}
