package bootstrap

import (
	"log"
	"os"
	"path/filepath"

	do "github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mytoy/config"
)

// NewContainer 初始化依赖注入容器，并注册全局单例依赖。
// 当前只注册了 *gorm.DB，后续可以在这里继续扩展。
func NewContainer() do.Injector {
	injector := do.New()

	// 注册配置管理器
	do.Provide(injector, func(i do.Injector) (*config.Manager, error) {
		return config.NewManager()
	})

	// 注册数据库依赖
	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {

		cfgMgr, err := do.Invoke[*config.Manager](i)
		if err != nil {
			return nil, err
		}
		snap := cfgMgr.Snapshot()
		dbPath := snap.Database.Path
		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, err
		}

		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		if sqlDB, err := db.DB(); err == nil && snap.Database.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(snap.Database.MaxOpenConns)
		}

		log.Printf("connected to SQLite database at %s", dbPath)

		return db, nil
	})

	return injector
}
