package initiallize

import (
	"github.com/gin-gonic/gin"
	"kubeimooc/middleware"
	"kubeimooc/router"
)

func Routers() *gin.Engine {
	r := gin.Default()
	// 跨域
	r.Use(middleware.Cors)

	//注册路由
	router.RouterGroupApp.ExampleRouterGroup.InitExample(r)
	router.RouterGroupApp.K8sRouterGroup.InitK8sRouter(r)

	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	return r
}
