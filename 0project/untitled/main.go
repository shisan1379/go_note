package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

func main() {
	router := gin.New()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello")
	})

	// 为http/2配置参数
	h2Handler := h2c.NewHandler(router, &http2.Server{})
	// 配置http服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2Handler,
	}
	// 启动http服务
	server.ListenAndServe()
}
