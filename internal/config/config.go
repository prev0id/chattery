package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppName       string
	HTTPAddress   string
	MessagesLimit int
}

func Init() *Config {
	config := &Config{}

	bindString(&config.AppName, "APP_NAME", "chattery")
	bindString(&config.HTTPAddress, "HTTP_ADDRESS", ":8080")
	bindInt(&config.MessagesLimit, "MESSAGES_LIMIT", 20)

	return config
}

func bindString(to *string, envName string, defaultValue string) {
	if fromEnv := os.Getenv(envName); fromEnv != "" {
		*to = fromEnv
		return
	}
	*to = defaultValue
}

func bindInt(to *int, envName string, defaultValue int) {
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		*to = defaultValue
		return
	}

	value, err := strconv.Atoi(fromEnv)
	if err != nil {
		*to = defaultValue
		return
	}

	*to = value
}
