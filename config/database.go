package config

import (
	"errors"
	"strings"

	"mytoy/util"
)

const databaseConfigPath = "config/database.toml"

type DatabaseConfig struct {
	Path         string `toml:"path"`
	MaxOpenConns int    `toml:"max_open_conns"`
}

func loadDatabaseConfig() (DatabaseConfig, error) {
	var cfg DatabaseConfig
	if err := util.ReadToml(databaseConfigPath, &cfg); err != nil {
		return cfg, err
	}
	if err := validateDatabase(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validateDatabase(cfg DatabaseConfig) error {
	if strings.TrimSpace(cfg.Path) == "" {
		return errors.New("database.path is required")
	}
	if cfg.MaxOpenConns <= 0 {
		return errors.New("database.max_open_conns is required and must be > 0")
	}
	return nil
}

func cloneDatabase(cfg DatabaseConfig) DatabaseConfig {
	return DatabaseConfig{
		Path:         cfg.Path,
		MaxOpenConns: cfg.MaxOpenConns,
	}
}
