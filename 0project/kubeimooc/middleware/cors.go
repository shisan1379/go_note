package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors(c *gin.Context) {
	// 设置允许访问的域名，这里使用 * 表示允许所有域名访问
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 设置允许的请求方法
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 设置允许的请求头
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	// 设置是否允许携带凭证（如 cookie）
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	// 如果是 OPTIONS 请求，通常是预检请求，直接返回 204 状态码
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	// 继续处理后续的请求
	c.Next()
}
