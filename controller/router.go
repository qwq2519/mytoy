package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	do "github.com/samber/do/v2"
)

// NewRouter 创建并返回 Gin 引擎。
// 这里接收 DI 容器，方便在注册路由时解析 service 等依赖。
func NewRouter(injector do.Injector) *gin.Engine {
	// 个人项目默认使用带日志和恢复的默认引擎
	r := gin.Default()

	// 健康检查 / 探活接口
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// TODO: 在这里继续注册各业务模块的路由，比如：
	// registerChatRoutes(r, injector)
	// registerCharacterRoutes(r, injector)

	return r
}
