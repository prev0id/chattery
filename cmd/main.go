package main

import (
	chatadapter "chattery/internal/adapter/postgres/chat"
	redisadapter "chattery/internal/adapter/redis"
	"chattery/internal/api"
	chatapi "chattery/internal/api/chat"
	"chattery/internal/client/redis"
	"chattery/internal/config"
	"chattery/internal/service/chat"
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
	chatDB := chatadapter.New(cfg, transactionManager)

	redisClient := redis.New(redisConn)
	redisAdapter := redisadapter.NewRedisAdapter(redisClient)

	_ = chat.New(chatDB, redisAdapter)

	server := api.
		NewServer(cfg).
		Register(
			chatapi.New(),
		)

	if err := server.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
