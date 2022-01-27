package config

import (
	"os"
)

var PG_HOST = getStringParameter("PG_HOST", "localhost")
var PG_PORT = getStringParameter("PG_PORT", "5444")
var POSTGRES_DB = getStringParameter("POSTGRES_DB", "scooterdb")
var POSTGRES_USER = getStringParameter("POSTGRES_USER", "scooteradmin")
var POSTGRES_PASSWORD = getStringParameter("POSTGRES_PASSWORD", "Megascooter!")
var GRPC_PORT = getStringParameter("GRPC_PORT", "9000")
var ORDER_GRPC_PORT = getStringParameter("ORDER_GRPC_PORT", "9999")
var KAFKA_BROKER = getStringParameter("KAFKA_BROKER", "localhost:9093")

func getStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}
