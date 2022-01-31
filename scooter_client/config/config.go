package config

import (
	"os"
)

var GRPC_PORT = getStringParameter("GRPC_PORT", "9000")
var ORDER_GRPC_PORT = getStringParameter("ORDER_GRPC_PORT", "9999")
var KAFKA_BROKER = getStringParameter("KAFKA_BROKER", "localhost:9093")
var SERVER_CONN_GRPC_ADDRESS = getStringParameter("SERVER_CONN_GRPC_ADDRESS", ":9000")

func getStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}
