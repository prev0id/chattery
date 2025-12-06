package config

import "os"

var (
	GRPCAddress = "localhost:50051"
	HTTPAddress = "localhost:8080"
)

func init() {
	initStringVar(&GRPCAddress, "GRPC_ADDRESS")
	initStringVar(&HTTPAddress, "HTTP_ADDRESS")
}

func initStringVar(value *string, key string) {
	if fromEnv, ok := os.LookupEnv(key); ok {
		*value = fromEnv
	}
}
