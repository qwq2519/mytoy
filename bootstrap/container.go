package bootstrap

import (
	"log"
	"os"
	"path/filepath"

	do "github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const defaultDBPath = "storage/db/app.sqlite"

// NewContainer 初始化依赖注入容器，并注册全局单例依赖。
// 当前只注册了 *gorm.DB，后续可以在这里继续扩展。
func NewContainer() do.Injector {
	injector := do.New()

	// 注册数据库依赖
	do.Provide(injector, func(i do.Injector) (*gorm.DB, error) {
		dbPath := os.Getenv("APP_DB_PATH")
		if dbPath == "" {
			dbPath = defaultDBPath
		}

		if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
			return nil, err
		}

		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		log.Printf("connected to SQLite database at %s", dbPath)

		return db, nil
	})

	return injector
}
