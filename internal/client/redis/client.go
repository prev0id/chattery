package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"chattery/internal/utils/errors"
)

type Client struct {
	conn *redis.Client
}

func New(client *redis.Client) *Client {
	return &Client{conn: client}
}

// func (c *Client) Get(ctx context.Context, key string) (string, error) {
// 	value, err := c.conn.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		return "", errors.E(err).Kind(errors.NotFound).Debug("key:" + key)
// 	}
// 	if err != nil {
// 		return "", errors.E(err).Debug("key:"+key, "c.conn.Get")
// 	}
// 	return value, nil
// }

func (c *Client) SetI64(ctx context.Context, key string, value int64, expiration time.Duration) error {
	if err := c.conn.Set(ctx, key, value, expiration).Err(); err != nil {
		return errors.E(err).Debug("key:"+key, "c.conn.Set")
	}
	return nil
}

func (c *Client) GetExI64(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	value, err := c.conn.GetEx(ctx, key, expiration).Int64()
	if err == redis.Nil {
		return 0, errors.E(err).Kind(errors.NotFound).Debug("key:" + key)
	}
	if err != nil {
		return 0, errors.E(err).Debug("key:"+key, "c.conn.GetEx")
	}
	return value, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.conn.Del(ctx, key).Err(); err != nil {
		return errors.E(err).Debug("key:"+key, "c.conn.Del")
	}
	return nil
}

func (c *Client) Publish(ctx context.Context, channel string, message string) error {
	if err := c.conn.Publish(ctx, channel, message).Err(); err != nil {
		return errors.E(err).Debug("channel:" + channel)
	}
	return nil
}

func (c *Client) Subscribe(ctx context.Context, channel string, sink chan<- string) {
	pubsub := c.conn.Subscribe(ctx, channel)
	defer pubsub.Close()

	for message := range pubsub.Channel() {
		sink <- message.Payload
	}

}
