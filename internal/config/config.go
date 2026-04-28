package config

import (
	"time"

	"chattery/internal/utils/bind"
)

type Config struct {
	Redis    Redis
	HTTP     HTTP
	Session  Session
	Postgres Postgres
	App      App
	Chat     Chat
	Cache    Cache
}

type App struct {
	Name    string
	Version string
	Debug   bool
}

type HTTP struct {
	Host string
	Port string
}

type Redis struct {
	Address  string
	Username string
	Password string
}

type Postgres struct {
	URL string
}

type Session struct {
	SecretKey  string
	Expiration time.Duration
}

type Chat struct {
	MessagesLimit int
}

type Cache struct {
	UserStoreSyncTimeout   time.Duration
	DMStoreSyncTimeout     time.Duration
	ServerStoreSyncTimeout time.Duration
}

func Init() *Config {
	return &Config{
		App: App{
			Name:    bind.EnvString("APP_NAME", "chattery"),
			Version: bind.EnvString("APP_VERSION", "local"),
			Debug:   bind.EnvBool("APP_DEBUG", false),
		},
		HTTP: HTTP{
			Host: bind.EnvString("HTTP_HOST", "localhost"),
			Port: bind.EnvString("HTTP_PORT", "8080"),
		},
		Session: Session{
			Expiration: bind.EnvDuration("SESSION_EXPIRATION", 5*time.Hour),
			SecretKey:  bind.EnvString("SESSION_KEY", "local-key"),
		},
		Redis: Redis{
			Address:  bind.EnvString("REDIS_ADDRESS", "localhost:6379"),
			Username: bind.EnvString("REDIS_USERNAME", "default"),
			Password: bind.EnvString("REDIS_PASSWORD", "redis_password"),
		},
		Chat: Chat{
			MessagesLimit: bind.EnvInt("MESSAGES_LIMIT", 20),
		},
		Cache: Cache{
			UserStoreSyncTimeout:   bind.EnvDuration("CACHE_USER_SYNC_TIMEOUT", 30*time.Second),
			ServerStoreSyncTimeout: bind.EnvDuration("CACHE_SERVER_SYNC_TIMEOUT", 30*time.Second),
			DMStoreSyncTimeout:     bind.EnvDuration("CACHE_DM_SYNC_TIMEOUT", 30*time.Second),
		},
		Postgres: Postgres{
			URL: bind.EnvString("POSTGRES_URL", "postgresql://user:password@localhost:5432/chattery?sslmode=disable"),
		},
	}
}
