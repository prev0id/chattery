package api

import (
	"chattery/internal/config"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type service interface {
	Register(*fiber.App)
}

type Server struct {
	app     *fiber.App
	address string
}

func NewServer(cfg *config.Config) *Server {
	server := &Server{
		address: cfg.HTTPAddress,
	}

	server.app = fiber.New(fiber.Config{
		Prefork:       true,
		ServerHeader:  cfg.AppName,
		CaseSensitive: true,
	})

	return server
}

func (s *Server) Register(services ...service) {
	for _, svc := range services {
		svc.Register(s.app)
	}
}

func (s *Server) Run() error {
	if err := s.app.Listen(s.address); err != nil {
		return fmt.Errorf("s.app.Listen: %w", err)
	}
	return nil
}
