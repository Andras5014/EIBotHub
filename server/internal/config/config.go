package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Config struct {
	Port       string
	DBPath     string
	StorageDir string
	AppSecret  string
	SeedDemo   bool
	GinMode    string
}

func Load() (Config, error) {
	cfg := defaultConfig()
	if err := loadConfigFile(&cfg); err != nil {
		return Config{}, err
	}
	if err := applyEnvOverrides(&cfg); err != nil {
		return Config{}, err
	}
	if err := validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func defaultConfig() Config {
	return Config{
		Port:       "8080",
		DBPath:     filepath.Join(".", "data", "opencommunity.db"),
		StorageDir: filepath.Join(".", "storage"),
		AppSecret:  "opencommunity-local-secret",
		SeedDemo:   true,
	}
}

func loadConfigFile(cfg *Config) error {
	path, ok, err := configFilePath()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return applyConfigFile(cfg, path)
}

func configFilePath() (string, bool, error) {
	if path := strings.TrimSpace(os.Getenv("CONFIG_FILE")); path != "" {
		if _, err := os.Stat(path); err != nil {
			return "", false, fmt.Errorf("read config file %q: %w", path, err)
		}
		return path, true, nil
	}

	for _, candidate := range configFileCandidates() {
		info, err := os.Stat(candidate)
		if err == nil {
			if info.IsDir() {
				continue
			}
			return candidate, true, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", false, fmt.Errorf("read config file %q: %w", candidate, err)
		}
	}
	return "", false, nil
}

func configFileCandidates() []string {
	candidates := []string{
		filepath.Join(".", "deploy", "config.json"),
		filepath.Join(".", "config.json"),
	}
	if exePath, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exePath)
		candidates = append(
			candidates,
			filepath.Join(exeDir, "deploy", "config.json"),
			filepath.Join(exeDir, "config.json"),
		)
	}

	result := make([]string, 0, len(candidates))
	seen := map[string]struct{}{}
	for _, candidate := range candidates {
		abs, err := filepath.Abs(candidate)
		if err != nil {
			abs = candidate
		}
		if _, ok := seen[abs]; ok {
			continue
		}
		seen[abs] = struct{}{}
		result = append(result, candidate)
	}
	return result
}

func applyConfigFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	raw := map[string]json.RawMessage{}
	if err := decoder.Decode(&raw); err != nil {
		return fmt.Errorf("parse config file %q: %w", path, err)
	}

	root := configPathRoot(path)
	aliases := configKeyAliases()
	seen := map[string]struct{}{}
	for key, value := range raw {
		canonical, ok := aliases[key]
		if !ok {
			return fmt.Errorf("unknown config key %q in %s", key, path)
		}
		if _, ok := seen[canonical]; ok {
			return fmt.Errorf("config file %s sets %q more than once", path, canonical)
		}
		seen[canonical] = struct{}{}

		if err := applyConfigValue(cfg, canonical, key, value, root); err != nil {
			return err
		}
	}
	return nil
}

func configKeyAliases() map[string]string {
	return map[string]string{
		"port":        "port",
		"app_port":    "port",
		"APP_PORT":    "port",
		"db_path":     "db_path",
		"dbPath":      "db_path",
		"DB_PATH":     "db_path",
		"storage_dir": "storage_dir",
		"storageDir":  "storage_dir",
		"STORAGE_DIR": "storage_dir",
		"app_secret":  "app_secret",
		"appSecret":   "app_secret",
		"APP_SECRET":  "app_secret",
		"seed_demo":   "seed_demo",
		"seedDemo":    "seed_demo",
		"SEED_DEMO":   "seed_demo",
		"gin_mode":    "gin_mode",
		"ginMode":     "gin_mode",
		"GIN_MODE":    "gin_mode",
	}
}

func applyConfigValue(cfg *Config, canonical, original string, value json.RawMessage, root string) error {
	switch canonical {
	case "port":
		port, err := decodeString(value, original)
		if err != nil {
			return err
		}
		cfg.Port = port
	case "db_path":
		dbPath, err := decodeString(value, original)
		if err != nil {
			return err
		}
		cfg.DBPath = resolveConfigPath(root, dbPath)
	case "storage_dir":
		storageDir, err := decodeString(value, original)
		if err != nil {
			return err
		}
		cfg.StorageDir = resolveConfigPath(root, storageDir)
	case "app_secret":
		appSecret, err := decodeString(value, original)
		if err != nil {
			return err
		}
		cfg.AppSecret = appSecret
	case "seed_demo":
		seedDemo, err := decodeBool(value, original)
		if err != nil {
			return err
		}
		cfg.SeedDemo = seedDemo
	case "gin_mode":
		ginMode, err := decodeString(value, original)
		if err != nil {
			return err
		}
		cfg.GinMode = ginMode
	}
	return nil
}

func configPathRoot(path string) string {
	dir := filepath.Dir(path)
	if filepath.Base(dir) == "deploy" {
		return filepath.Dir(dir)
	}
	return dir
}

func resolveConfigPath(root, value string) string {
	if value == "" || filepath.IsAbs(value) {
		return value
	}
	return filepath.Clean(filepath.Join(root, value))
}

func decodeString(raw json.RawMessage, key string) (string, error) {
	var value string
	if err := json.Unmarshal(raw, &value); err == nil {
		return value, nil
	}
	var number json.Number
	if err := json.Unmarshal(raw, &number); err == nil {
		return number.String(), nil
	}
	return "", fmt.Errorf("config key %q must be a string", key)
}

func decodeBool(raw json.RawMessage, key string) (bool, error) {
	var value bool
	if err := json.Unmarshal(raw, &value); err == nil {
		return value, nil
	}
	var text string
	if err := json.Unmarshal(raw, &text); err == nil {
		parsed, err := strconv.ParseBool(strings.TrimSpace(text))
		if err != nil {
			return false, fmt.Errorf("config key %q must be a bool: %w", key, err)
		}
		return parsed, nil
	}
	return false, fmt.Errorf("config key %q must be a bool", key)
}

func applyEnvOverrides(cfg *Config) error {
	if value := os.Getenv("APP_PORT"); value != "" {
		cfg.Port = value
	}
	if value := os.Getenv("DB_PATH"); value != "" {
		cfg.DBPath = value
	}
	if value := os.Getenv("STORAGE_DIR"); value != "" {
		cfg.StorageDir = value
	}
	if value := os.Getenv("APP_SECRET"); value != "" {
		cfg.AppSecret = value
	}
	if value := os.Getenv("SEED_DEMO"); value != "" {
		seedDemo, err := strconv.ParseBool(strings.TrimSpace(value))
		if err != nil {
			return fmt.Errorf("parse SEED_DEMO: %w", err)
		}
		cfg.SeedDemo = seedDemo
	}
	if value := os.Getenv("GIN_MODE"); value != "" {
		cfg.GinMode = value
	}
	return nil
}

func validate(cfg Config) error {
	if strings.TrimSpace(cfg.Port) == "" {
		return errors.New("APP_PORT cannot be empty")
	}
	if strings.TrimSpace(cfg.DBPath) == "" {
		return errors.New("DB_PATH cannot be empty")
	}
	if strings.TrimSpace(cfg.StorageDir) == "" {
		return errors.New("STORAGE_DIR cannot be empty")
	}
	if strings.TrimSpace(cfg.AppSecret) == "" {
		return errors.New("APP_SECRET cannot be empty")
	}
	if cfg.GinMode != "" {
		switch cfg.GinMode {
		case "debug", "release", "test":
		default:
			return fmt.Errorf("GIN_MODE must be debug, release, or test, got %q", cfg.GinMode)
		}
	}
	return nil
}
