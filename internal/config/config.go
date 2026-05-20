package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ServerAddr  string
	DBDriver    string
	DBDSN       string
	CORSOrigin  string
	AuthSecret  string
	UploadDir   string
	Environment string
}

func MustLoad() Config {
	if err := loadDotEnv(".env"); err != nil {
		panic(fmt.Errorf("load .env: %w", err))
	}

	cfg := Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DBDriver:    getEnv("DB_DRIVER", "sqlite"),
		DBDSN:       getEnv("DB_DSN", "file:app.db?cache=shared"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "*"),
		AuthSecret:  getEnv("AUTH_SECRET", "dev-change-me"),
		UploadDir:   getEnv("UPLOAD_DIR", filepath.FromSlash("static/uploads")),
		Environment: strings.ToLower(getEnv("APP_ENV", "development")),
	}

	mustNotBlank("DB_DRIVER", cfg.DBDriver)
	mustNotBlank("DB_DSN", cfg.DBDSN)
	mustNotBlank("AUTH_SECRET", cfg.AuthSecret)
	mustNotBlank("UPLOAD_DIR", cfg.UploadDir)

	return cfg
}

func mustNotBlank(name, value string) {
	if strings.TrimSpace(value) == "" {
		panic(fmt.Sprintf("%s is required", name))
	}
}

func getEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func loadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, `"'`)
		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}

	return scanner.Err()
}
