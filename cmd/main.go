package main

import (
	"context"

	chatadapter "chattery/internal/adapter/postgres/chat"
	useradapter "chattery/internal/adapter/postgres/user"
	redisadapter "chattery/internal/adapter/redis"
	"chattery/internal/api"
	signalingapi "chattery/internal/api/signaling"
	userapi "chattery/internal/api/user"
	"chattery/internal/client/redis"
	"chattery/internal/config"
	"chattery/internal/service/chat"
	"chattery/internal/service/user"
	"chattery/internal/utils/database"
	"chattery/internal/utils/logger"
	"chattery/internal/utils/transaction"
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
	userDB := useradapter.New(transactionManager)

	redisClient := redis.New(redisConn)
	redisAdapter := redisadapter.NewRedisAdapter(redisClient)

	_ = chat.New(chatDB, redisAdapter)
	userService := user.New(userDB, redisAdapter, transactionManager)

	server := api.
		NewServer(cfg).
		Register(
			signalingapi.New(),
			userapi.New(userService),
		)

	if err := server.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
