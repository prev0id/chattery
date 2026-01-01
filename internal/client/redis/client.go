package redisclient

import (
	"chattery/internal/config"

	"github.com/gomodule/redigo/redis"
)

type Client struct {
	pool *redis.Pool
}

func New(cfg config.Config) *Client {
	pool := &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisAddress)
		},
	}
	return &Client{pool: pool}
}

func (c *Client) Conn() redis.Conn {
	return c.pool.Get()
}
