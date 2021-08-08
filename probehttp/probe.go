package probehttp

import (
	"fmt"
	"net/http"
	"simple-http-probe/probe"

	"github.com/gin-gonic/gin"
)

// 定义HttpProbe方法
func HttpProbe(c *gin.Context) {

	// 解析参数
	host := c.Query("host")
	isHttps := c.Query("is_https")

	// 校验入参
	if host == "" {
		c.String(http.StatusBadRequest, "empty host")
	}
	// 默认http协议
	schema := "http"
	if isHttps == "1" {
		schema = "https"
	}
	url := fmt.Sprintf("%s://%s", schema, host)
	res := probe.DoHttpProbe(url)

	c.String(http.StatusOK, res)
}
