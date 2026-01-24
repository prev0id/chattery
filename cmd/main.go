package main

import (
	"context"

	chat_adapter "chattery/internal/adapter/postgres/chat"
	useradapter "chattery/internal/adapter/postgres/user"
	redisadapter "chattery/internal/adapter/redis"
	"chattery/internal/api"
	chat_api "chattery/internal/api/chat"
	signalingapi "chattery/internal/api/signaling"
	user_api "chattery/internal/api/user"
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
	chatDB := chat_adapter.New(cfg, transactionManager)
	userDB := useradapter.New(transactionManager)

	redisClient := redis.New(redisConn)
	redisAdapter := redisadapter.NewRedisAdapter(redisClient)

	chatService := chat.New(chatDB, redisAdapter, transactionManager)
	userService := user.New(userDB, redisAdapter, transactionManager)

	server := api.
		NewServer(cfg).
		Register(
			signalingapi.New(),
			user_api.New(userService),
			chat_api.New(userService, chatService),
			// web_api.New(),
		)

	if err := server.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
