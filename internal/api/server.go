package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	"chattery/internal/config"
	"chattery/internal/domain"
)

const MiB = 1024 * 1024

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 30 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 120 * time.Second
	shutdownTimeout   = 10 * time.Second
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
		address: net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port),
	}
	server.mux = chi.NewRouter()
	server.mux.Use(
		middleware.RequestSize(cfg.Image.MaxFileSize+MiB),
		httplog.RequestLogger(slog.Default(), nil),
		middleware.RequestID,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)
	server.mux.NotFound(redirectUnknownPath)

	return server
}

func (s *Server) UseMiddleware(middlewares ...func(http.Handler) http.Handler) *Server {
	s.mux.Use(middlewares...)
	return s
}

func (s *Server) Register(services ...service) *Server {
	for _, svc := range services {
		s.mux.Route(svc.Pattern(), svc.Route)
	}
	return s
}

func (s *Server) Run(ctx context.Context) error {
	slog.Info("starting server", slog.String("address", s.address))

	server := &http.Server{
		Addr:              s.address,
		Handler:           s.mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	errCh := make(chan error)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}
}

func redirectUnknownPath(w http.ResponseWriter, r *http.Request) {
	if domain.UserIDFromContext(r.Context()) != domain.UserIsUnknown {
		http.Redirect(w, r, "/app/dm", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
