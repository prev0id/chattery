package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"chattery/internal/config"
	"chattery/internal/utils/errors"
)

func PostgresConnection(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.Postgres.URL)
	if err != nil {
		return nil, errors.E(err).Debug("pgx.Connect")
	}
	if err := conn.Ping(ctx); err != nil {
		return nil, errors.E(err).Debug("conn.Ping")
	}
	return conn, nil
}

func RedisConnection(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:       cfg.Redis.Address,
		Username:   cfg.Redis.Username,
		Password:   cfg.Redis.Password,
		ClientName: cfg.App.Name,
		Protocol:   3,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.E(err).Debug("client.Ping")
	}
	return client, nil
}
