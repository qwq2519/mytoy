package config

import (
	"errors"
	"strings"

	"mytoy/util"
)

const loggingConfigPath = "config/logging.toml"

type LoggingConfig struct {
	Level         string `toml:"level"`
	File          string `toml:"file"`
	EnableConsole bool   `toml:"enable_console"`
}

func loadLoggingConfig() (LoggingConfig, error) {
	var cfg LoggingConfig
	if err := util.ReadToml(loggingConfigPath, &cfg); err != nil {
		return cfg, err
	}
	if err := validateLogging(cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func validateLogging(cfg LoggingConfig) error {
	if strings.TrimSpace(cfg.Level) == "" {
		return errors.New("logging.level is required")
	}
	if strings.TrimSpace(cfg.File) == "" {
		return errors.New("logging.file is required")
	}
	// EnableConsole 是布尔类型，无法区分缺失与 false，这里允许为 false。
	return nil
}

func cloneLogging(cfg LoggingConfig) LoggingConfig {
	return LoggingConfig{
		Level:         cfg.Level,
		File:          cfg.File,
		EnableConsole: cfg.EnableConsole,
	}
}
