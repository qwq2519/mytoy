package config

import (
	"sync/atomic"
)

// AppConfig 聚合各模块配置。
type AppConfig struct {
	Server   ServerConfig   `toml:"server"`
	Database DatabaseConfig `toml:"database"`
	Logging  LoggingConfig  `toml:"logging"`
}

// Manager 提供线程安全的配置快照读写。
type Manager struct {
	snapshot atomic.Pointer[AppConfig]
}

func NewManager() (*Manager, error) {
	m := &Manager{}
	if _, err := m.Reload(); err != nil {
		return nil, err
	}
	return m, nil
}

// Snapshot 返回当前配置副本，避免外部修改内部状态。
func (m *Manager) Snapshot() AppConfig {
	if snap := m.snapshot.Load(); snap != nil {
		return cloneAppConfig(*snap)
	}
	return AppConfig{}
}

// Reload 从文件加载所有模块配置，若文件缺失或字段缺失则返回错误。
func (m *Manager) Reload() (AppConfig, error) {
	cfg, err := loadAll()
	if err != nil {
		return AppConfig{}, err
	}
	m.snapshot.Store(&cfg)
	return cloneAppConfig(cfg), nil
}

func loadAll() (AppConfig, error) {
	server, err := loadServerConfig()
	if err != nil {
		return AppConfig{}, err
	}
	database, err := loadDatabaseConfig()
	if err != nil {
		return AppConfig{}, err
	}
	logging, err := loadLoggingConfig()
	if err != nil {
		return AppConfig{}, err
	}
	return AppConfig{
		Server:   server,
		Database: database,
		Logging:  logging,
	}, nil
}

func cloneAppConfig(cfg AppConfig) AppConfig {
	return AppConfig{
		Server:   cloneServer(cfg.Server),
		Database: cloneDatabase(cfg.Database),
		Logging:  cloneLogging(cfg.Logging),
	}
}
