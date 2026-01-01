package main

import (
	"chattery/internal/api"
	chatapi "chattery/internal/api/chat"
	"chattery/internal/config"
	"chattery/internal/utils/errors"
	"chattery/internal/utils/logger"
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
		errors.E(err).Debug("server.Run").LogFatal()
	}
}
