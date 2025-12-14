package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	do "github.com/samber/do/v2"

	"mytoy/bootstrap"
	"mytoy/config"
	"mytoy/controller"
)

func main() {
	// 初始化依赖注入容器（包括数据库等）
	injector := bootstrap.NewContainer()

	cfgManager, err := do.Invoke[*config.Manager](injector)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	snapshot := cfgManager.Snapshot()

	if snapshot.Server.Mode == "" {
		log.Fatalf("server.mode is required")
	}
	gin.SetMode(snapshot.Server.Mode)

	// 初始化路由
	router := controller.NewRouter(injector)

	if snapshot.Server.Port <= 0 {
		log.Fatalf("server.port is required and must be > 0")
	}
	addr := fmt.Sprintf(":%d", snapshot.Server.Port)

	log.Printf("starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
