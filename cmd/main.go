package main

import (
	"context"
	"os/signal"
	"syscall"

	dm_adapter "chattery/internal/adapter/postgres/dm"
	server_adapter "chattery/internal/adapter/postgres/server"
	user_adapter "chattery/internal/adapter/postgres/user"
	redis_adapter "chattery/internal/adapter/redis"
	s3_adapter "chattery/internal/adapter/s3"
	"chattery/internal/api"
	dm_api "chattery/internal/api/dm"
	image_api "chattery/internal/api/image"
	server_api "chattery/internal/api/server"
	user_api "chattery/internal/api/user"
	web_api "chattery/internal/api/web"
	websocket_api "chattery/internal/api/websocket"
	"chattery/internal/client/redis"
	"chattery/internal/client/s3"
	"chattery/internal/config"
	"chattery/internal/service/dm"
	image_service "chattery/internal/service/image"
	"chattery/internal/service/server"
	"chattery/internal/service/text_topic"
	"chattery/internal/service/user"
	"chattery/internal/service/voice_topic"
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
	appCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
	s3Client, err := s3.New(cfg.S3)
	if err != nil {
		logger.Fatal(err, "s3.New")
	}

	transactionManager := transaction.NewManager(postgresConn)
	dmDB := dm_adapter.New(transactionManager)
	serverDB := server_adapter.New(transactionManager)
	userDB := user_adapter.New(transactionManager)

	userStore := user_store.New(userDB)
	dmStore := dm_store.New(dmDB)
	serverStore := server_store.New(serverDB)

	redisClient := redis.New(redisConn)
	redisAdapter := redis_adapter.NewRedisAdapter(redisClient, cfg.Session.SecretKey)
	s3Adapter := s3_adapter.New(s3Client, cfg)

	dmService := dm.New(dmDB, transactionManager, redisAdapter, userStore, dmStore, cfg)
	serverService := server.New(serverDB, transactionManager, serverStore, cfg)
	textTopicService := text_topic.New(serverDB, transactionManager, redisAdapter, userStore, cfg)
	voiceTopicService, err := voice_topic.New(redisAdapter, serverService, userStore, cfg)
	if err != nil {
		logger.Fatal(err, "voice_topic.New")
	}
	userService := user.New(userDB, redisAdapter, transactionManager, cfg)
	imageService, err := image_service.New(userDB, s3Adapter, cfg)
	if err != nil {
		logger.Fatal(err, "image_service.New")
	}

	if err := syncer.Start(cfg.Cache.UserStoreSyncTimeout, userStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}
	if err := syncer.Start(cfg.Cache.DMStoreSyncTimeout, dmStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}
	if err := syncer.Start(cfg.Cache.ServerStoreSyncTimeout, serverStore); err != nil {
		logger.Fatal(err, "syncer.New")
	}
	voiceTopicService.Start(appCtx)

	wsManagerInstance := ws_manager.New(redisAdapter, dmService, serverService, voiceTopicService)

	apiServer := api.
		NewServer(cfg).
		UseMiddleware(userService.SessionMiddleware).
		Register(
			websocket_api.New(userService, wsManagerInstance),
			user_api.New(userService),
			image_api.New(userService, userStore, imageService, cfg),
			dm_api.New(userService, dmService, userStore),
			web_api.New(),
			server_api.New(userService, serverService, textTopicService, userStore),
		)

	if err := apiServer.Run(appCtx); err != nil {
		logger.Fatal(err, "server.Run")
	}
}
