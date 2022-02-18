package config

import (
	"os"
)

var HTTP_PORT = GetStringParameter("HTTP_PORT", "8085")
var PG_HOST = GetStringParameter("PG_HOST", "localhost")
var PG_PORT = GetStringParameter("PG_PORT", "5444")
var POSTGRES_DB = GetStringParameter("POSTGRES_DB", "scooterdb")
var POSTGRES_USER = GetStringParameter("POSTGRES_USER", "scooteradmin")
var POSTGRES_PASSWORD = GetStringParameter("POSTGRES_PASSWORD", "Megascooter!")
var GRPC_PORT = GetStringParameter("GRPC_PORT", "9000")
var ORDER_GRPC_PORT = GetStringParameter("ORDER_GRPC_PORT", "9999")
var MONO_TEMPLATES_PATH = GetStringParameter("MONO_TEMPLATES_PATH", "../scooter_server/templates/")
var KAFKA_BROKER = GetStringParameter("KAFKA_BROKER", "localhost:9093")

func GetStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}
