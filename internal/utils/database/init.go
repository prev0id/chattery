package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"chattery/internal/config"
	"chattery/internal/utils/errutil"
)

func PostgresConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Postgres.URL)
	if err != nil {
		return nil, errutil.E(err).Debug("pgx.Connect")
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, errutil.E(err).Debug("conn.Ping")
	}
	return pool, nil
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
		return nil, errutil.E(err).Debug("client.Ping")
	}
	return client, nil
}

func IsConstraintViolation(err error, constraintName string) bool {
	if err == nil {
		return false
	}

	pgxErr := &pgconn.PgError{}
	if ok := errors.As(err, &pgxErr); ok {
		return pgxErr.ConstraintName == constraintName
	}
	return false
}

func NotFound(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, pgx.ErrNoRows)
}
