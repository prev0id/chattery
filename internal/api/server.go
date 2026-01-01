package api

import (
	"chattery/internal/config"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

type service interface {
	Pattern() string
	Route(chi.Router)
}

type Server struct {
	mux     *chi.Mux
	address string
}

func NewServer(cfg *config.Config) *Server {
	server := &Server{
		address: cfg.HTTPAddress,
	}

	server.mux = chi.NewRouter()
	server.mux.Use(
		httplog.RequestLogger(slog.Default(), nil),
		middleware.RequestID,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Heartbeat("/ping"),
	)

	return server
}

func (s *Server) Register(services ...service) {
	for _, svc := range services {
		s.mux.Route(svc.Pattern(), svc.Route)
	}
}

func (s *Server) Run() error {
	slog.Info("starting server", slog.String("address", s.address))

	if err := http.ListenAndServe(s.address, s.mux); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}
