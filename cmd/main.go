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
	user_api "chattery/internal/api/user"
	web_api "chattery/internal/api/web"
	websocket_api "chattery/internal/api/websocket"
	"chattery/internal/client/redis"
	"chattery/internal/config"
	"chattery/internal/service/dm"
	"chattery/internal/service/server"
	"chattery/internal/service/text_topic"
	"chattery/internal/service/user"
	ws_manager "chattery/internal/service/websocket_manager"
	dm_store "chattery/internal/store/dm"
	server_store "chattery/internal/store/server"
	"chattery/internal/store/syncer"
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

	userStore := user_store.New(userDB)
	dmStore := dm_store.New(dmDB)
	serverStore := server_store.New(serverDB)

	redisClient := redis.New(redisConn)
	redisAdapter := redis_adapter.NewRedisAdapter(redisClient)

	dmService := dm.New(dmDB, transactionManager, redisAdapter, userStore, dmStore, cfg)
	serverService := server.New(serverDB, transactionManager, serverStore, cfg)
	textTopicService := text_topic.New(serverDB, transactionManager, redisAdapter, userStore, cfg)
	userService := user.New(userDB, redisAdapter, transactionManager, cfg)

	if err := syncer.Start(cfg.Cache.UserStoreSyncTimeout, userStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}
	if err := syncer.Start(cfg.Cache.DMStoreSyncTimeout, dmStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}
	if err := syncer.Start(cfg.Cache.ServerStoreSyncTimeout, serverStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}

	wsManagerInstance := ws_manager.New(redisAdapter, dmService, serverService)

	apiServer := api.
		NewServer(cfg).
		UseMiddleware(userService.SessionMiddleware).
		Register(
			websocket_api.New(userService, wsManagerInstance),
			user_api.New(userService),
			image_api.New(userStore),
			dm_api.New(userService, dmService, userStore),
			web_api.New(),
			server_api.New(userService, serverService, textTopicService, userStore),
		)

	if err := apiServer.Run(); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
