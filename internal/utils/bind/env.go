package bind

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

func EnvBool(envName string, defaultValue bool) bool {
	fromEnv := os.Getenv(envName)
	if fromEnv == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(fromEnv)
	if err != nil {
		return defaultValue
	}

	return value
}

func PathParamI64[T ~int64](r *http.Request, paramName string) (T, error) {
	rawParam := chi.URLParam(r, paramName)
	param, err := strconv.ParseInt(rawParam, 10, 64)
	return T(param), err
}
