package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5"

	postgres_adapter "chattery/backend/user_service/internal/adapter/postgres"
	"chattery/backend/user_service/internal/api"
	postgres_client "chattery/backend/user_service/internal/client/postgres"
	"chattery/backend/user_service/internal/config"
	"chattery/backend/user_service/internal/service/user"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	conn, err := pgx.Connect(ctx, config.DBString)
	if err != nil {
		log.Fatalf("pgx.Connect: %s", err.Error())
	}

	pgClient := postgres_client.New(conn)
	pgAdapter := postgres_adapter.New(pgClient)
	userService := user.New(pgAdapter)

	if err := api.StartUserServer(ctx, userService); err != nil {
		log.Fatalf("api.StartUserServer: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")
	cancel()
}
