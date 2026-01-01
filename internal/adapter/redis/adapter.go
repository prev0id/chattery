package redisadapter

import (
	"chattery/internal/utils/errors"

	"github.com/gomodule/redigo/redis"
)

type client interface {
	Get() redis.Conn
}

type Adapter struct {
	client client
}

func NewRedisAdapter(client client) *Adapter {
	return &Adapter{client: client}
}

func (r *Adapter) Publish(channel string, msg []byte) error {
	conn := r.client.Get()
	defer conn.Close()

	if _, err := conn.Do("PUBLISH", channel, msg); err != nil {
		return errors.E(err).Debug("conn.Do PUBLISH")
	}
	return nil
}

func (r *Adapter) Subscribe(channel string, sink chan []byte) {
	conn := r.client.Get()
	psc := redis.PubSubConn{Conn: conn}

	_ = psc.Subscribe(channel)

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			sink <- v.Data
		}
	}
}
