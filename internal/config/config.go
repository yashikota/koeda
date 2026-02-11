package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	TTL time.Duration `json:"ttl"`
}

func GetConfigDir() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config")
	}
	return filepath.Join(configDir, "koeda"), nil
}

func Load() (*Config, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "config.json")

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{TTL: 24 * time.Hour}, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	if cfg.TTL == 0 {
		cfg.TTL = 24 * time.Hour
	}
	return &cfg, nil
}

func Save(cfg *Config) error {
	dir, err := GetConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	path := filepath.Join(dir, "config.json")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cfg)
}
