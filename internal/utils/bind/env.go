package bind

import (
	"os"
	"strconv"
	"time"
)

func EnvString(envName string, defaultValue string) string {
	if fromEnv := os.Getenv(envName); fromEnv != "" {
		return fromEnv
	}

	return defaultValue
}

func EnvInt(envName string, defaultValue int) int {
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(fromEnv)
	if err != nil {
		return defaultValue
	}

	return value
}

func EnvDuration(envName string, defaultValue time.Duration) time.Duration {
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(fromEnv)
	if err != nil {
		return defaultValue
	}

	return value
}
