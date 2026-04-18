package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	Port       string
	DBPath     string
	StorageDir string
	AppSecret  string
	SeedDemo   bool
}

func Load() Config {
	return Config{
		Port:       envOrDefault("APP_PORT", "8080"),
		DBPath:     envOrDefault("DB_PATH", filepath.Join(".", "data", "opencommunity.db")),
		StorageDir: envOrDefault("STORAGE_DIR", filepath.Join(".", "storage")),
		AppSecret:  envOrDefault("APP_SECRET", "opencommunity-local-secret"),
		SeedDemo:   envOrDefault("SEED_DEMO", "true") != "false",
	}
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
