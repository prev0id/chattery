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
		slog.String("app.name", cfg.App.Name),
		slog.String("app.version", cfg.App.Version),
	)
	slog.SetDefault(logger)
}
