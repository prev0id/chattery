package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"

	"chattery/internal/config"
	"chattery/internal/domain"
	"chattery/internal/utils/errors"
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

func Error(err error, message string, attr ...slog.Attr) {
	e := errors.E(err)
	attr = append(attr,
		slog.Group("error",
			slog.String("kind", e.GetKind().String()),
			slog.Any("debug", e.GetDebug()),
			slog.String("message", e.GetMessage()),
			slog.Any("err", e.GetError()),
		),
	)

	slog.LogAttrs(context.Background(), slog.LevelError, message, attr...)
}

func ErrorCtx(ctx context.Context, err error, message string, attr ...slog.Attr) {
	requestID := middleware.GetReqID(ctx)
	userID := domain.UserIDFromContext(ctx)

	Error(err,
		"request ended with an error",
		slog.String("request_id", requestID),
		slog.Int64("user_id", userID.I64()),
	)
}

func Fatal(err error, message string, attr ...slog.Attr) {
	Error(err, message, attr...)
	os.Exit(1)
}
