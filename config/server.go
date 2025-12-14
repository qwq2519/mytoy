package config

import (
	"errors"
	"strings"

	"mytoy/util"
)

const serverConfigPath = "config/server.toml"

type ServerConfig struct {
	Port int    `toml:"port"`
	Mode string `toml:"mode"`
}

func loadServerConfig() (ServerConfig, error) {
	var cfg ServerConfig
	if err := util.ReadToml(serverConfigPath, &cfg); err != nil {
		return cfg, err
	}
	if err := validateServer(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validateServer(cfg ServerConfig) error {
	if cfg.Port <= 0 {
		return errors.New("server.port is required and must be > 0")
	}
	if strings.TrimSpace(cfg.Mode) == "" {
		return errors.New("server.mode is required")
	}
	return nil
}

func cloneServer(cfg ServerConfig) ServerConfig {
	return ServerConfig{
		Port: cfg.Port,
		Mode: cfg.Mode,
	}
}
