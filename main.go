package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	do "github.com/samber/do/v2"

	"mytoy/bootstrap"
	"mytoy/config"
	"mytoy/controller"
)

func main() {
	// 初始化依赖注入容器（包括数据库等）
	injector := bootstrap.NewContainer()

	logger, err := do.Invoke[*slog.Logger](injector)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}

	cfgManager, err := do.Invoke[*config.Manager](injector)
	if err != nil {
		logger.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	snapshot := cfgManager.Snapshot()

	if snapshot.Server.Mode == "" {
		logger.Error("server.mode is required")
		os.Exit(1)
	}
	gin.SetMode(snapshot.Server.Mode)

	// 初始化路由
	router := controller.NewRouter(injector)

	if snapshot.Server.Port <= 0 {
		logger.Error("server.port is required and must be > 0", "port", snapshot.Server.Port)
		os.Exit(1)
	}
	addr := fmt.Sprintf(":%d", snapshot.Server.Port)

	logger.Info("starting server", "addr", addr)

	if err := router.Run(addr); err != nil {
		logger.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}
