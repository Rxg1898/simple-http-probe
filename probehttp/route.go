package probehttp

import (
	"simple-http-probe/config"

	"github.com/gin-gonic/gin"
)

func Setup(c *config.Config) *gin.Engine {
	// 初始化gin
	r := gin.Default()
	// 绑定路由
	Routers(r)
	// 启动
	return r
}

// 路由
func Routers(r *gin.Engine) {
	api := r.Group("/api")
	api.GET("/probe/http", HttpProbe)
}
