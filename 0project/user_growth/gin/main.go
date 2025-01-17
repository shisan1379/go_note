package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/testdata/protoexample"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()

	// 此路由将匹配 /user/xxx 但不会匹配 /user/ 或 /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	//此handler 将匹配 /user/xxx 和 /user/xxx/send
	// 如果没有其他路由匹配 /user/xxx 它将重定向到上面的 /user/:name
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/login", func(c *gin.Context) {})
		v1.POST("/submit", func(c *gin.Context) {})
		v1.POST("/read", func(c *gin.Context) {})

	}
	// 全局中间件
	// Logger 中间件将日志写入 gin.DefaultWriter 即使将 GIN_MODE 设置为 release
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())
	// Recover 中间件会 捕获任何 panic 如果发生 panic 的话会写入 500
	router.Use(gin.Recovery())
	// 认证路由组
	group := router.Group("/", func(context *gin.Context) {
		//执行认证逻辑
	})

	// 认证后才能执行
	{
		group.GET("/a", func(context *gin.Context) {})
		group.GET("/a", func(context *gin.Context) {})
		group.GET("/a", func(context *gin.Context) {})
	}

	router.GET("/json", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"message": "hello world", "status": http.StatusOK})
	})
	router.GET("/xml", func(context *gin.Context) {
		context.XML(http.StatusOK, gin.H{"message": "hello world", "status": http.StatusOK})
	})
	router.GET("/yaml", func(context *gin.Context) {
		context.YAML(http.StatusOK, gin.H{"message": "hello world", "status": http.StatusOK})
	})

	router.GET("/protoBuf", func(context *gin.Context) {

		reps := []int64{1, 2}
		lable := "test"
		data := &protoexample.Test{
			Label: &lable,
			Reps:  reps,
		}

		// 数据会在响应中变为二进制数据
		// 将输出被 protoexample.Test Protobuf 序列化了的数据
		context.ProtoBuf(http.StatusOK, data)
	})

	router.GET("/html", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "hello world",
		})
	})

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server.ListenAndServe()
}
func cookie() {
	// 创建一个默认的 Gin 引擎
	r := gin.Default()

	// 定义一个路由用于获取 Cookie
	r.GET("/getcookie", func(c *gin.Context) {
		// 尝试从请求中获取名为 "my_cookie" 的 Cookie
		cookie, err := c.Cookie("my_cookie")
		if err != nil {
			// 如果找不到 Cookie，则返回错误响应
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cookie not found"})
			return
		}
		// 如果找到了 Cookie，则返回其值
		c.JSON(http.StatusOK, gin.H{"cookie_value": cookie})
	})

	// 定义一个路由用于设置 Cookie
	r.GET("/setcookie", func(c *gin.Context) {
		// 创建一个新的 Cookie 对象
		cookie := &http.Cookie{
			Name:     "my_cookie",  // Cookie 的名称
			Value:    "some_value", // Cookie 的值
			MaxAge:   60,           // Cookie 的有效期（秒）
			Path:     "/",          // Cookie 的路径
			Domain:   "",           // Cookie 的域名（留空表示当前域名）
			Secure:   false,        // 是否仅通过 HTTPS 发送 Cookie（false 表示不限制）
			HttpOnly: true,         // 是否仅通过 HTTP/HTTPS 协议访问（true 表示不能通过 JavaScript 访问）
		}
		// 将 Cookie 添加到响应中
		http.SetCookie(c.Writer, cookie)
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "Cookie set successfully"})
	})

	// 定义一个路由用于删除 Cookie
	r.GET("/deletecookie", func(c *gin.Context) {
		// 创建一个用于删除的 Cookie 对象（通过设置 MaxAge 为 -1）
		cookie := &http.Cookie{
			Name:     "my_cookie",
			Value:    "", // 值可以留空，因为 MaxAge 为 -1 会导致浏览器删除这个 Cookie
			MaxAge:   -1, // 负值表示删除 Cookie
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
		}
		// 将删除 Cookie 的指令添加到响应中
		http.SetCookie(c.Writer, cookie)
		// 返回成功响应
		c.JSON(http.StatusOK, gin.H{"message": "Cookie deleted successfully"})
	})

	// 启动 Gin 引擎并监听端口 8080
	r.Run(":8080")
}
