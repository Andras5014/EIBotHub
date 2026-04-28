package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsDeployConfig(t *testing.T) {
	tempDir := t.TempDir()
	deployDir := filepath.Join(tempDir, "deploy")
	if err := os.MkdirAll(deployDir, 0o755); err != nil {
		t.Fatalf("create deploy dir: %v", err)
	}
	configContent := []byte(`{
  "port": "9090",
  "db_path": "./data/test.db",
  "storage_dir": "./files",
  "app_secret": "file-secret",
  "seed_demo": false,
  "gin_mode": "release"
}`)
	if err := os.WriteFile(filepath.Join(deployDir, "config.json"), configContent, 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	for _, key := range []string{"CONFIG_FILE", "APP_PORT", "DB_PATH", "STORAGE_DIR", "APP_SECRET", "SEED_DEMO", "GIN_MODE"} {
		t.Setenv(key, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Port != "9090" {
		t.Fatalf("Port = %q, want 9090", cfg.Port)
	}
	if cfg.DBPath != filepath.Join("data", "test.db") {
		t.Fatalf("DBPath = %q, want %q", cfg.DBPath, filepath.Join("data", "test.db"))
	}
	if cfg.StorageDir != "files" {
		t.Fatalf("StorageDir = %q, want files", cfg.StorageDir)
	}
	if cfg.AppSecret != "file-secret" {
		t.Fatalf("AppSecret = %q, want file-secret", cfg.AppSecret)
	}
	if cfg.SeedDemo {
		t.Fatal("SeedDemo = true, want false")
	}
	if cfg.GinMode != "release" {
		t.Fatalf("GinMode = %q, want release", cfg.GinMode)
	}
}

func TestLoadEnvOverridesConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")
	configContent := []byte(`{
  "port": "9090",
  "db_path": "./data/test.db",
  "storage_dir": "./files",
  "app_secret": "file-secret",
  "seed_demo": false,
  "gin_mode": "release"
}`)
	if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	envDBPath := filepath.Join(tempDir, "env.db")
	t.Setenv("CONFIG_FILE", configPath)
	t.Setenv("APP_PORT", "7070")
	t.Setenv("DB_PATH", envDBPath)
	t.Setenv("STORAGE_DIR", "env-storage")
	t.Setenv("APP_SECRET", "env-secret")
	t.Setenv("SEED_DEMO", "true")
	t.Setenv("GIN_MODE", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Port != "7070" {
		t.Fatalf("Port = %q, want 7070", cfg.Port)
	}
	if cfg.DBPath != envDBPath {
		t.Fatalf("DBPath = %q, want %q", cfg.DBPath, envDBPath)
	}
	if cfg.StorageDir != "env-storage" {
		t.Fatalf("StorageDir = %q, want env-storage", cfg.StorageDir)
	}
	if cfg.AppSecret != "env-secret" {
		t.Fatalf("AppSecret = %q, want env-secret", cfg.AppSecret)
	}
	if !cfg.SeedDemo {
		t.Fatal("SeedDemo = false, want true")
	}
	if cfg.GinMode != "debug" {
		t.Fatalf("GinMode = %q, want debug", cfg.GinMode)
	}
}
