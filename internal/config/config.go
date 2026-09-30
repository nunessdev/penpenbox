package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// Config to save API key and Steam ID.
type Config struct {
	SteamID string `json:"steam_id"`
	APIKey  string `json:"api_key"`
}

// Check if both values exist
func (c Config) HasSteam() bool {
	return c.SteamID != "" && c.APIKey != ""
}

// Get path to config
func path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "penpenbox", "config.json"), nil
}

// Load reads the config file. Missing file: TUI asks for configuration
func Load() (Config, error) {
	var cfg Config

	p, err := path()
	if err != nil {
		return cfg, err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cfg, nil // first run: empty config, no error
		}
		return cfg, err // a real problem
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// Save writes the config to disk.
func Save(cfg Config) error {
	p, err := path()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(p, data, 0o600)
}