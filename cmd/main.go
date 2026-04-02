package main

import (
	"context"

	dm_adapter "chattery/internal/adapter/postgres/dm"
	server_adapter "chattery/internal/adapter/postgres/server"
	user_adapter "chattery/internal/adapter/postgres/user"
	redis_adapter "chattery/internal/adapter/redis"
	"chattery/internal/api"
	dm_api "chattery/internal/api/dm"
	image_api "chattery/internal/api/image"
	server_api "chattery/internal/api/server"
	signaling_api "chattery/internal/api/signaling"
	user_api "chattery/internal/api/user"
	web_api "chattery/internal/api/web"
	"chattery/internal/client/redis"
	"chattery/internal/config"
	hub_pkg "chattery/internal/hub"
	"chattery/internal/service/dm"
	"chattery/internal/service/server"
	"chattery/internal/service/user"
	syncer "chattery/internal/store/syncer"
	user_store "chattery/internal/store/user"
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
	dmDB := dm_adapter.New(transactionManager)
	serverDB := server_adapter.New(transactionManager)
	userDB := user_adapter.New(transactionManager)

	redisClient := redis.New(redisConn)
	redisAdapter := redis_adapter.NewRedisAdapter(redisClient)

	dmService := dm.New(dmDB, transactionManager, redisAdapter, cfg)
	serverService := server.New(serverDB, transactionManager, redisAdapter, cfg)
	userService := user.New(userDB, redisAdapter, transactionManager, cfg)

	userStore := user_store.New(userDB)

	if err := syncer.Start(cfg.Cache.UserStoreSyncTimeout, userStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}

	hubInstance := hub_pkg.New(redisAdapter, dmService, serverService, userStore)

	server := api.
		NewServer(cfg).
		UseMiddleware(userService.SessionMiddleware).
		Register(
			signaling_api.New(userService, hubInstance),
			user_api.New(userService),
			image_api.New(userStore),
			dm_api.New(userService, dmService, userStore),
			web_api.New(),
			server_api.New(userService, serverService, userStore),
		)

	if err := server.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
