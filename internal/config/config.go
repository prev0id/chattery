package config

import (
	"os"
	"time"

	"chattery/internal/utils/bind"
)

type Config struct {
	Redis    Redis
	Session  Session
	HTTP     HTTP
	Postgres Postgres
	App      App
	Voice    Voice
	Cache    Cache
	Chat     Chat
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
	SecretKey     string
	CookieDomain  string
	Expiration    time.Duration
	RefreshBefore time.Duration
}

type Chat struct {
	MessagesLimit int
}

type Cache struct {
	UserStoreSyncTimeout   time.Duration
	DMStoreSyncTimeout     time.Duration
	ServerStoreSyncTimeout time.Duration
}

type Voice struct {
	NodeID     string
	STUNServer string
	OwnerTTL   time.Duration
}

func Init() *Config {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "local"
	}

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
			Expiration:    bind.EnvDuration("SESSION_EXPIRATION", 5*time.Hour),
			RefreshBefore: bind.EnvDuration("SESSION_REFRESH_BEFORE", 30*time.Minute),
			SecretKey:     bind.EnvString("SESSION_KEY", "local-key"),
			CookieDomain:  bind.EnvString("SESSION_COOKIE_DOMAIN", ""),
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
		Voice: Voice{
			NodeID:     bind.EnvString("VOICE_NODE_ID", hostname),
			OwnerTTL:   bind.EnvDuration("VOICE_OWNER_TTL", 30*time.Second),
			STUNServer: bind.EnvString("VOICE_STUN_SERVER", "stun:stun.l.google.com:19302"),
		},
		Postgres: Postgres{
			URL: bind.EnvString("POSTGRES_URL", "postgresql://user:password@localhost:5432/chattery?sslmode=disable"),
		},
	}
}
