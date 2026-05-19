package config

import (
	"log"
	"os"
)

type Config struct {
	ServerAddr string
	DBDriver   string
	DBDSN      string
	CORSOrigin string
}

func MustLoad() Config {
	cfg := Config{
		ServerAddr: getEnv("SERVER_ADDR", ":8080"),
		DBDriver:   getEnv("DB_DRIVER", "sqlite3"),
		DBDSN:      getEnv("DB_DSN", ":memory:"),
		CORSOrigin: getEnv("CORS_ORIGIN", "*"),
	}

	if cfg.DBDriver == "" || cfg.DBDSN == "" {
		log.Fatal("DB_DRIVER and DB_DSN are required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
