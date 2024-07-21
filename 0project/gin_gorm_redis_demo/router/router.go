package router

import (
	"gin_gorm_redis_demo/config"
	"gin_gorm_redis_demo/controllers"
	logger "gin_gorm_redis_demo/pkg/log"
	"github.com/gin-contrib/sessions"
	sessions_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()

	r.Use(gin.LoggerWithConfig(logger.LoggerToFile()))
	r.Use(logger.Recover)

	store, _ := sessions_redis.NewStore(10, "tcp", config.RedisAddress, "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello world")
	})

	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", controllers.UserController{}.Register)
		userGroup.POST("/login", controllers.UserController{}.Login)
	}
	playerGroup := r.Group("/player")
	{
		playerGroup.POST("/getPlayers", controllers.PlayerController{}.GetPlayers)
	}
	return r
}
