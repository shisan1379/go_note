package controllers

import (
	"github.com/gin-gonic/gin"
	logger "go_project_demo/pkg/log"
	"strconv"
)

type User struct {
}

func (u *User) GetUserInfo(c *gin.Context) {
	//获取url路径中的参数
	//例如： /user/123

	num1 := 1
	num2 := 0
	num3 := num1 / num2

	id := c.Param("id")
	logger.Write(id, "user")
	ReturnSuccess(c, 0, "success", "user info:"+id+strconv.Itoa(num3), 10)
}

func (u *User) GetUserInfo2(c *gin.Context) {
	//获取url中?后的参数
	id := c.Query("id")
	ReturnSuccess(c, 0, "success", "user info:"+id, 10)
}

func (u *User) GetList(c *gin.Context) {

	ReturnError(c, 40004, "error")
}
