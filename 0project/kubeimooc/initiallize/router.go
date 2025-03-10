package initiallize

import (
	"github.com/gin-gonic/gin"
	"kubeimooc/router"
)

func Routers() {
	r := gin.Default()

	router.RouterGroupApp.ExampleRouterGroup.InitExample(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
