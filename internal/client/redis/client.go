package redis

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"chattery/internal/utils/errutil"
)

const inf = "+inf"

type Client struct {
	conn *redis.Client
}

func New(client *redis.Client) *Client {
	return &Client{conn: client}
}

// func (c *Client) Get(ctx context.Context, key string) (string, error) {
// 	value, err := c.conn.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		return "", errutil.E(err).Kind(errutil.NotFound).Debug("key:" + key)
// 	}
// 	if err != nil {
// 		return "", errutil.E(err).Debug("key:"+key, "c.conn.Get")
// 	}
// 	return value, nil
// }

func (c *Client) SetI64(ctx context.Context, key string, value int64, expiration time.Duration) error {
	if err := c.conn.Set(ctx, key, value, expiration).Err(); err != nil {
		return errutil.E(err).Debug("key:"+key, "c.conn.Set")
	}
	return nil
}

func (c *Client) GetExI64(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	value, err := c.conn.GetEx(ctx, key, expiration).Int64()
	if errors.Is(err, redis.Nil) {
		return 0, errutil.E(err).Kind(errutil.NotFound).Debug("key:" + key)
	}
	if err != nil {
		return 0, errutil.E(err).Debug("key:"+key, "c.conn.GetEx")
	}
	return value, nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	if err := c.conn.Del(ctx, key).Err(); err != nil {
		return errutil.E(err).Debug("key:"+key, "c.conn.Del")
	}
	return nil
}

func (c *Client) ZAddI64(ctx context.Context, key string, value int64) error {
	member := redis.Z{
		Member: value,
		Score:  float64(time.Now().Unix()),
	}
	if err := c.conn.ZAdd(ctx, key, member).Err(); err != nil {
		return errutil.E(err).Debug("key:"+key, "c.conn.SAdd")
	}
	return nil
}

func (c *Client) ZMembersI64(ctx context.Context, key string, threshold time.Duration) ([]int64, error) {
	rangeOpts := &redis.ZRangeBy{
		Min: strconv.FormatInt(time.Now().Add(-threshold).Unix(), 10),
		Max: inf,
	}

	members, err := c.conn.ZRangeByScore(ctx, key, rangeOpts).Result()
	if err != nil {
		return nil, errutil.E(err).Debug("key:"+key, "c.conn.SMembers")
	}

	result := make([]int64, 0, len(members))
	for _, member := range members {
		i64, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			return nil, errutil.E(err).Debug("key:"+key, "can't parse int64", member)
		}
		result = append(result, i64)
	}
	return result, nil
}

func (c *Client) Publish(ctx context.Context, channel string, message string) error {
	if err := c.conn.Publish(ctx, channel, message).Err(); err != nil {
		return errutil.E(err).Debug("channel:" + channel)
	}
	return nil
}

func (c *Client) Subscribe(ctx context.Context, channel string, sink chan<- string, done <-chan struct{}) {
	pubsub := c.conn.Subscribe(ctx, channel)
	defer pubsub.Close()

	for {
		select {
		case <-done:
			close(sink)
			return
		case message, ok := <-pubsub.Channel():
			if !ok {
				return
			}
			sink <- message.Payload
		}
	}
}

func (c *Client) ZRemoveI64(ctx context.Context, key string, value int64) error {
	if err := c.conn.ZRem(ctx, key, value).Err(); err != nil {
		return errutil.E(err).Debug("key:"+key, "c.conn.ZRem")
	}
	return nil
}
