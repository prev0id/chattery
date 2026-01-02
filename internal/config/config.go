package config

import (
	"chattery/internal/utils/bind"
	"time"
)

type Config struct {
	App     App
	Redis   Redis
	Session Session
	Http    Http
	Chat    Chat
}

type App struct {
	Name    string
	Version string
}

type Http struct {
	Host string
	Port string
}

type Redis struct {
	Address  string
	Username string
	Password string
}

type Session struct {
	Expiration time.Duration
	SecretKey  string
}

type Chat struct {
	MessagesLimit int
}

func Init() *Config {
	return &Config{
		App: App{
			Name:    bind.EnvString("APP_NAME", "chattery"),
			Version: bind.EnvString("APP_VERSION", "local"),
		},
		Http: Http{
			Host: bind.EnvString("HTTP_HOST", "localhost"),
			Port: bind.EnvString("HTTP_PORT", ":8080"),
		},
		Session: Session{
			Expiration: bind.EnvDuration("SESSION_EXPIRATION", 5*time.Minute),
			SecretKey:  bind.EnvString("SESSION_KEY", "local-key"),
		},
		Redis: Redis{
			Address:  bind.EnvString("REDIS_ADDRESS", "localhost:6379"),
			Username: bind.EnvString("REDIS_USERNAME", "redis_user"),
			Password: bind.EnvString("REDIS_PASSWORD", "redis_password"),
		},
		Chat: Chat{
			MessagesLimit: bind.EnvInt("MESSAGES_LIMIT", 20),
		},
	}
}
