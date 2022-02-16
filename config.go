package main

import (
	"os"
)

var HTTP_PORT = getStringParameter("HTTP_PORT", "8085")
var PG_HOST = getStringParameter("PG_HOST", "localhost")
var PG_PORT = getStringParameter("PG_PORT", "5444")
var POSTGRES_DB = getStringParameter("POSTGRES_DB", "scooterdb")
var POSTGRES_USER = getStringParameter("POSTGRES_USER", "scooteradmin")
var POSTGRES_PASSWORD = getStringParameter("POSTGRES_PASSWORD", "Megascooter!")
var GRPC_PORT = getStringParameter("GRPC_PORT", "9000")
var ORDER_GRPC_PORT = getStringParameter("ORDER_GRPC_PORT", "9999")
var MONO_TEMPLATES_PATH = getStringParameter("MONO_TEMPLATES_PATH", "../scooter_server/templates/")
var KAFKA_BROKER = getStringParameter("KAFKA_BROKER", "localhost:9093")

func getStringParameter(paramName, defaultValue string) string {
	result, ok := os.LookupEnv(paramName)
	if !ok {
		result = defaultValue
	}
	return result
}
