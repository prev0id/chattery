package api

import (
	"chattery/internal/config"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

const MiB = 1024 * 1024

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
		address: net.JoinHostPort(cfg.Http.Host, cfg.Http.Port),
	}
	server.mux = chi.NewRouter()
	server.mux.Use(
		middleware.RequestSize(2*MiB),
		httplog.RequestLogger(slog.Default(), nil),
		middleware.RequestID,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Heartbeat("/ping"),
	)

	return server
}

func (s *Server) Register(services ...service) *Server {
	for _, svc := range services {
		s.mux.Route(svc.Pattern(), svc.Route)
	}
	return s
}

func (s *Server) Run() error {
	slog.Info("starting server", slog.String("address", s.address))

	if err := http.ListenAndServe(s.address, s.mux); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	return nil
}
