package bootstrap

import (
	"io"
	"log"
	"log/slog"
	"mytoy/util/logging"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	do "github.com/samber/do/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"mytoy/config"
)

// loggingResource 用于在 DI Shutdown 阶段关闭日志文件句柄。
type loggingResource struct {
	closer io.Closer
}

func (r *loggingResource) Shutdown() error {
	if r == nil || r.closer == nil {
		return nil
	}
	return r.closer.Close()
}

// dbResource 用于在 DI Shutdown 阶段关闭底层 sql.DB。
type dbResource struct {
	db *gorm.DB
}

func (r *dbResource) Shutdown() error {
	if r == nil || r.db == nil {
		return nil
	}
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// NewContainer 初始化依赖注入容器，并注册全局单例依赖。
func NewContainer() do.Injector {
	injector := do.New()

	// 注册配置管理器
	do.Provide(injector, func(i do.Injector) (*config.Manager, error) {
		return config.NewManager()
	})

	// 注册日志依赖
	do.Provide(injector, func(i do.Injector) (*slog.Logger, error) {
		cfgMgr, err := do.Invoke[*config.Manager](i)
		if err != nil {
			return nil, err
		}
		snap := cfgMgr.Snapshot()

		logger, writer, closer, err := logging.NewLogger(snap.Logging)
		if err != nil {
			return nil, err
		}
		do.ProvideValue(i, &loggingResource{closer: closer})

		slog.SetDefault(logger)

		stdLog := slog.NewLogLogger(logger.Handler(), slog.LevelInfo)
		stdLog.SetFlags(0)
		log.SetFlags(0)
		log.SetOutput(stdLog.Writer())

		gin.DefaultWriter = writer
		gin.DefaultErrorWriter = writer

		return logger, nil
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

		logger, err := do.Invoke[*slog.Logger](i)
		if err != nil {
			return nil, err
		}
		logger.Info("connected to SQLite database", "path", dbPath)

		do.ProvideValue(i, &dbResource{db: db})

		return db, nil
	})

	return injector
}
