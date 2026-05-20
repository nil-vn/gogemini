package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	ServerAddr           string
	DBDriver             string
	DBDSN                string
	DBMaxOpen            int
	DBMaxIdle            int
	DBMaxLifeMs          int
	CORSOrigin           string
	AuthSecret           string
	UploadDir            string
	Environment          string
	DefaultAdminUsername string
	DefaultAdminPassword string
	DefaultAdminEmail    string
	DefaultAdminRole     string
	DefaultAdminStatus   string
}

func MustLoad() Config {
	if err := loadDotEnv(".env"); err != nil {
		panic(fmt.Errorf("load .env: %w", err))
	}

	cfg := Config{
		ServerAddr:           getEnv("SERVER_ADDR", ":8080"),
		DBDriver:             getEnv("DB_DRIVER", "sqlite"),
		DBDSN:                getEnv("DB_DSN", "file:app.db?cache=shared"),
		DBMaxOpen:            getEnvInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdle:            getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBMaxLifeMs:          getEnvInt("DB_CONN_MAX_LIFETIME_MS", 300000),
		CORSOrigin:           getEnv("CORS_ORIGIN", "*"),
		AuthSecret:           getEnv("AUTH_SECRET", "dev-change-me"),
		UploadDir:            getEnv("UPLOAD_DIR", filepath.FromSlash("static/uploads")),
		Environment:          strings.ToLower(getEnv("APP_ENV", "development")),
		DefaultAdminUsername: getEnv("DEFAULT_ADMIN_USERNAME", ""),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", ""),
		DefaultAdminEmail:    getEnv("DEFAULT_ADMIN_EMAIL", ""),
		DefaultAdminRole:     getEnv("DEFAULT_ADMIN_ROLE", "admin"),
		DefaultAdminStatus:   getEnv("DEFAULT_ADMIN_STATUS", "active"),
	}

	mustNotBlank("DB_DRIVER", cfg.DBDriver)
	mustNotBlank("DB_DSN", cfg.DBDSN)
	mustNotBlank("AUTH_SECRET", cfg.AuthSecret)
	mustNotBlank("UPLOAD_DIR", cfg.UploadDir)
	mustPositive("DB_MAX_OPEN_CONNS", cfg.DBMaxOpen)
	mustPositive("DB_MAX_IDLE_CONNS", cfg.DBMaxIdle)
	mustNonNegative("DB_CONN_MAX_LIFETIME_MS", cfg.DBMaxLifeMs)

	if cfg.DBMaxIdle > cfg.DBMaxOpen {
		panic("DB_MAX_IDLE_CONNS must be <= DB_MAX_OPEN_CONNS")
	}

	return cfg
}

func mustPositive(name string, value int) {
	if value <= 0 {
		panic(fmt.Sprintf("%s must be > 0", name))
	}
}

func mustNonNegative(name string, value int) {
	if value < 0 {
		panic(fmt.Sprintf("%s must be >= 0", name))
	}
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

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		panic(fmt.Sprintf("%s must be an integer", key))
	}
	return value
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
