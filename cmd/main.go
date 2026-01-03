package main

import (
	chatadapter "chattery/internal/adapter/postgres/chat"
	"chattery/internal/api"
	chatapi "chattery/internal/api/chat"
	"chattery/internal/client/redis"
	"chattery/internal/config"
	"chattery/internal/utils/database"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/transaction"
	"context"
)

func main() {
	appCtx := context.Background()

	cfg := config.Init()

	logger.Init(cfg)

	postgresConn, err := database.PostgresConnection(appCtx, cfg)
	if err != nil {
		logger.Fatal(err, "database.PostgresConnection")
	}
	redisConn, err := database.RedisConnection(appCtx, cfg)
	if err != nil {
		logger.Fatal(err, "database.RedisConnection")
	}

	transactionManager := transaction.NewManager(postgresConn)
	_ = redis.New(redisConn)

	_ = chatadapter.New(cfg, transactionManager)

	server := api.
		NewServer(cfg).
		Register(
			chatapi.New(),
		)

	if err := server.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
