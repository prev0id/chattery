package config

import (
	"log"
	"os"
	"time"
)

var (
	GRPCAddress = "localhost:50051"
	HTTPAddress = "localhost:8080"
	DBString    = "postgresql://user:password@localhost:5432/chattery?sslmode=disable"
	Expiration  = time.Hour * 24 * 7
	SigningKey  = ""
)

func init() {
	initStringVar(&GRPCAddress, "GRPC_ADDRESS")
	initStringVar(&HTTPAddress, "HTTP_ADDRESS")
	initStringVar(&DBString, "DBSTRING")
	initDuration(&Expiration, "EXPIRATION")
	initStringVar(&SigningKey, "SIGNING_KEY")
}

func initStringVar(value *string, key string) {
	if fromEnv, ok := os.LookupEnv(key); ok {
		*value = fromEnv
	}
}

func initDuration(value *time.Duration, key string) {
	fromEnv, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	newDuration, err := time.ParseDuration(fromEnv)
	if err != nil {
		log.Fatalf("time.ParseDuration key=%q: %s", key, err.Error())
	}
	*value = newDuration
}
