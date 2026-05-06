package config

import (
	"os"
	"time"

	"chattery/internal/utils/bind"
)

type Config struct {
	Redis    Redis
	HTTP     HTTP
	Postgres Postgres
	S3       S3
	App      App
	Session  Session
	Voice    Voice
	Image    Image
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

type S3 struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	UseSSL          bool
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

type Image struct {
	MaxFileSize int64
	MaxWidth    int
	MaxHeight   int
	JPEGQuality int
}

type Cache struct {
	UserStoreSyncTimeout   time.Duration
	DMStoreSyncTimeout     time.Duration
	ServerStoreSyncTimeout time.Duration
}

type Voice struct {
	NodeID        string
	STUNServer    string
	ICEPublicIP   string
	ICEUDPPortMin int
	ICEUDPPortMax int
	OwnerTTL      time.Duration
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
		Image: Image{
			MaxFileSize: bind.EnvInt64("IMAGE_MAX_FILE_SIZE", 2*1024*1024),
			MaxWidth:    bind.EnvInt("IMAGE_MAX_WIDTH", 1920),
			MaxHeight:   bind.EnvInt("IMAGE_MAX_HEIGHT", 1080),
			JPEGQuality: bind.EnvInt("IMAGE_JPEG_QUALITY", 75),
		},
		Cache: Cache{
			UserStoreSyncTimeout:   bind.EnvDuration("CACHE_USER_SYNC_TIMEOUT", 30*time.Second),
			ServerStoreSyncTimeout: bind.EnvDuration("CACHE_SERVER_SYNC_TIMEOUT", 30*time.Second),
			DMStoreSyncTimeout:     bind.EnvDuration("CACHE_DM_SYNC_TIMEOUT", 30*time.Second),
		},
		Voice: Voice{
			NodeID:        bind.EnvString("VOICE_NODE_ID", hostname),
			OwnerTTL:      bind.EnvDuration("VOICE_OWNER_TTL", 30*time.Second),
			STUNServer:    bind.EnvString("VOICE_STUN_SERVER", "stun:stun.l.google.com:19302"),
			ICEPublicIP:   bind.EnvString("VOICE_ICE_PUBLIC_IP", ""),
			ICEUDPPortMin: bind.EnvInt("VOICE_ICE_UDP_PORT_MIN", 0),
			ICEUDPPortMax: bind.EnvInt("VOICE_ICE_UDP_PORT_MAX", 0),
		},
		Postgres: Postgres{
			URL: bind.EnvString("POSTGRES_URL", "postgresql://user:password@localhost:5432/chattery?sslmode=disable"),
		},
		S3: S3{
			Endpoint:        bind.EnvString("S3_ENDPOINT", "localhost:9000"),
			AccessKeyID:     bind.EnvString("S3_ACCESS_KEY_ID", "user"),
			SecretAccessKey: bind.EnvString("S3_SECRET_ACCESS_KEY", "password"),
			Bucket:          bind.EnvString("S3_BUCKET", "chattery"),
			UseSSL:          bind.EnvBool("S3_USE_SSL", false),
		},
	}
}
