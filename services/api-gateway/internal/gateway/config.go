package gateway

import (
	"io"
	"os"
)

type Config struct {
	ListenAddr string
	JWTSecret  string
	Routes     map[string]string
	LogOutput  io.Writer
}

func NewConfig() Config {
	return Config{
		ListenAddr: envOrDefault("GATEWAY_ADDR", ":8080"),
		JWTSecret:  envOrDefault("GATEWAY_JWT_SECRET", "secret"),
		LogOutput:  os.Stdout,
		Routes: map[string]string{
			"/auth":          envOrDefault("AUTH_SERVICE_URL", "http://localhost:8081"),
			"/users":         envOrDefault("USER_SERVICE_URL", "http://localhost:8082"),
			"/products":      envOrDefault("CATALOG_SERVICE_URL", "http://localhost:8083"),
			"/cart":          envOrDefault("CART_SERVICE_URL", "http://localhost:8084"),
			"/orders":        envOrDefault("ORDER_SERVICE_URL", "http://localhost:8085"),
			"/payments":      envOrDefault("PAYMENT_SERVICE_URL", "http://localhost:8086"),
			"/notifications": envOrDefault("NOTIFICATION_SERVICE_URL", "http://localhost:8087"),
		},
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
