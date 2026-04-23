package domain

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	URL   string `toml:"url"`
	Token string `toml:"token"`
}

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", path, err)
	}

	if cfg.URL == "" {
		return nil, fmt.Errorf("missing 'url' in %s", path)
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("missing 'token' in %s", path)
	}

	return &cfg, nil
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".config", "paperclip", "config.toml"), nil
}

func EnsureConfigPermissions() error {
	path, err := configPath()
	if err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil // file doesn't exist yet, nothing to fix
	}
	if info.Mode().Perm() != 0600 {
		return os.Chmod(path, 0600)
	}
	return nil
}
