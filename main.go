package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"mytoy/bootstrap"
	"mytoy/controller"
)

func main() {
	// 如果未显式设置 GIN_MODE，则默认使用 release 模式，避免生产日志过于冗余。
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化依赖注入容器（包括数据库等）
	injector := bootstrap.NewContainer()

	// 初始化路由
	router := controller.NewRouter(injector)

	// 端口可以从环境变量 PORT 读取，默认使用 8080
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
