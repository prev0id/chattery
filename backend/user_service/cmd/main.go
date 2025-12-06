package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"chattery/backend/user_service/internal/api"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	if err := api.StartChatServer(ctx); err != nil {
		log.Fatalf("api.StartChatServer: %s", err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")
	cancel()
}
