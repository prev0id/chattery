package config

import (
	"chattery/internal/utils/bind"
	"time"
)

type Config struct {
	AppName           string
	AppVersion        string
	HTTPAddress       string
	RedisAddress      string
	MessagesLimit     int
	SessionExpiration time.Duration
	SessionSecretKey  string
}

func Init() *Config {
	return &Config{
		AppName:           bind.EnvString("APP_NAME", "chattery"),
		AppVersion:        bind.EnvString("APP_VERSION", ""),
		HTTPAddress:       bind.EnvString("HTTP_ADDRESS", ":8080"),
		RedisAddress:      bind.EnvString("REDIS_ADDRESS", ""),
		MessagesLimit:     bind.EnvInt("MESSAGES_LIMIT", 20),
		SessionExpiration: bind.EnvDuration("SESSION_EXPIRATION", 24*time.Hour),
		SessionSecretKey:  bind.EnvString("SESSION_KEY", ""),
	}
}
