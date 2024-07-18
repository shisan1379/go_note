package router

import (
	"github.com/gin-gonic/gin"
	"go_project_demo/controllers"
	logger "go_project_demo/pkg/log"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	r.Use(logger.Recover)

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})

	userRouter := r.Group("/user")
	{
		user := controllers.User{}

		userRouter.GET("/info/:id", user.GetUserInfo)
		userRouter.GET("/add/:name", user.AddUser)
		userRouter.GET("/update/:id/:name", user.UpdateUser)
		userRouter.GET("/info2", user.GetUserInfo2)
		userRouter.POST("/list", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "user list")
		})
		userRouter.PUT("/add", user.GetList)
		userRouter.DELETE("/delete", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "user delete")
		})
	}
	orderGroup := r.Group("/order")
	{
		order := controllers.Order{}
		orderGroup.POST("/info", order.GetList)
		orderGroup.POST("/info2", order.GetList2)
		orderGroup.POST("/search", order.SearchList)
	}
	return r
}
