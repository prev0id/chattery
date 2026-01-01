package main

import (
	"chattery/internal/api"
	chatapi "chattery/internal/api/chat"
	"chattery/internal/config"
	"chattery/internal/utils/logger"
	"log"
)

func main() {
	cfg := config.Init()

	logger.Init(cfg)

	server := api.NewServer(cfg)

	chatApi := chatapi.New()

	server.Register(
		chatApi,
	)

	if err := server.Run(); err != nil {
		log.Fatalf("server.Run: %s", err.Error())
	}
}
