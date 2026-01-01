package logger

import (
	"chattery/internal/config"
	"log/slog"
	"os"

	"github.com/go-chi/httplog/v3"
)

func Init(cfg *config.Config) {
	logFormat := httplog.SchemaECS.Concise(true)

	logger := slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			ReplaceAttr: logFormat.ReplaceAttr,
		},
	)).With(
		slog.String("app", cfg.AppName),
		slog.String("version", cfg.AppVersion),
	)
	slog.SetDefault(logger)
}
