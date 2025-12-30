package main

import (
	"chattery/internal/api"
	"chattery/internal/config"

	"github.com/gofiber/fiber/v2/log"
)

func main() {
	cfg := config.Init()

	server := api.NewServer(cfg)

	if err := server.Run(); err != nil {
		log.Fatalf("server.Run: %s", err.Error())
	}
}
