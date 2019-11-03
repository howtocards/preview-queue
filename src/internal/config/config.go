package config

import (
	"os"
	"strconv"

	"github.com/go-openapi/loads"
	"github.com/howtocards/preview-queue/src/internal/api/generated/restapi"
)

// Log field names.
const (
	LogRemote     = "remote" // aligned IPv4:Port "   192.168.0.42:1234 "
	LogFunc       = "func"   // RPC method name, REST resource path
	LogHTTPMethod = "httpMethod"
	LogHTTPStatus = "httpStatus"
)

// Default values.
var (
	oapiHost, oapiPort, oapiBasePath = apiConfig()

	APIHost     = strGetenv("API_HOST", oapiHost)
	APIPort     = intGetenv("API_PORT", oapiPort)
	APIBasePath = strGetenv("API_BASEPATH", oapiBasePath)
)

func intGetenv(name string, def int) int {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

func strGetenv(name, def string) string {
	value := os.Getenv(name)
	if value == "" {
		return def
	}
	return value
}

func apiConfig() (host string, port int, basePath string) {
	var err error
	port = 8080
	host, err = os.Hostname()
	if err != nil {
		host = "localhost"
	}

	spec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		return "", 0, ""
	}

	return host, port, spec.BasePath()
}
