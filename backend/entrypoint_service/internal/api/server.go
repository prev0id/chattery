package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler interface {
	GetPathPrefix() string
	RegisterRoutes() chi.Router
}

type Server struct {
	router   chi.Router
	addresss string
}

func NewServer(ctx context.Context, handlers ...Handler) *Server {
	router := chi.NewRouter()

	router.Use(middleware.Logger) // should be first
	router.Use(middleware.RedirectSlashes)
	router.Use(middleware.Recoverer) // should be last

	for _, handler := range handlers {
		router.Mount(handler.GetPathPrefix(), handler.RegisterRoutes())
	}

	return &Server{
		router: router,
	}
}

func (s *Server) ListenAndServer() error {
	slog.Info("staring server", slog.String("address", s.addresss))
	if err := http.ListenAndServe(s.addresss, s.router); err != nil {
		return fmt.Errorf("http.ListenAndServe: %w", err)
	}
	slog.Info("shutting down server")
	return nil
}
