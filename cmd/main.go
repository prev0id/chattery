package main

import (
	"context"

	chat_adapter "chattery/internal/adapter/postgres/chat"
	user_adapter "chattery/internal/adapter/postgres/user"
	redis_adapter "chattery/internal/adapter/redis"
	"chattery/internal/api"
	chat_api "chattery/internal/api/chat"
	signaling_api "chattery/internal/api/signaling"
	user_api "chattery/internal/api/user"
	web_api "chattery/internal/api/web"
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
	userDB := user_adapter.New(transactionManager)

	redisClient := redis.New(redisConn)
	redisAdapter := redis_adapter.NewRedisAdapter(redisClient)

	chatService := chat.New(chatDB, redisAdapter, transactionManager)
	userService := user.New(userDB, redisAdapter, transactionManager)

	server := api.
		NewServer(cfg).
		Register(
			signaling_api.New(chatService),
			user_api.New(userService),
			chat_api.New(userService, chatService),
			web_api.New(),
		)

	if err := server.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
