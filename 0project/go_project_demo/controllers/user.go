package controllers

import (
	"github.com/gin-gonic/gin"
	"go_project_demo/dao/dao"
	"go_project_demo/models"
	logger "go_project_demo/pkg/log"
	"strconv"
)

type User struct {
}

func (u *User) GetUserInfo(c *gin.Context) {
	//获取url路径中的参数
	//例如： /user/123

	//num1 := 1
	//num2 := 0
	//num3 := num1 / num2

	idStr := c.Param("id")
	logger.Write(idStr, "user")

	id, _ := strconv.Atoi(idStr)
	test, err := models.GetUserTest(id)
	if err != nil {
		logger.Error(map[string]interface{}{"err": err})
	}
	ReturnSuccess(c, 0, "success", test, 10)
}

func (u *User) GetUserInfo2(c *gin.Context) {
	//获取url中?后的参数
	id := c.Query("id")
	ReturnSuccess(c, 0, "success", "user info:"+id, 10)
}

func (u *User) GetList(c *gin.Context) {

	ReturnError(c, 40004, "error")
}

func (u *User) AddUser(c *gin.Context) {
	name := c.Param("name")

	user, _ := models.AddUser(name)

	ReturnSuccess(c, 0, "success", user, 10)
}
func (u *User) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	name := c.Param("name")
	id, _ := strconv.Atoi(idStr)

	//更新所有列
	var user models.User
	dao.Db.Where("id = ?", id).First(&user)
	user.Name = name
	dao.Db.Save(&user)

	//更新单个列
	dao.Db.Model(&models.User{}).Where("id = ?", id).Update("name", "hello")

	ReturnSuccess(c, 0, "success", user, 10)
}
func (u *User) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	dao.Db.Delete(&models.User{}, id)
	ReturnSuccess(c, 0, "success", nil, 10)
}
func (u *User) FindUser(c *gin.Context) {
	name := c.Param("name")
	var users []models.User
	dao.Db.Where("name = ?", name).Find(&users)
	ReturnSuccess(c, 0, "success", users, 10)
}
