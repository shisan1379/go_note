package controllers

import (
	"gin_gorm_redis_demo/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserController struct{}
type UserApi struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (r UserController) Register(c *gin.Context) {
	name := c.DefaultPostForm("name", "")
	pwd := c.DefaultPostForm("pwd", "")
	confirmPwd := c.DefaultPostForm("confirmPwd", "")

	if name == "" || pwd == "" || confirmPwd == "" {
		ReturnError(c, 4001, "请输入正确的信息")
		return
	}
	if pwd != confirmPwd {
		ReturnError(c, 4001, "两次输入密码不相同")
		return
	}

	user, err := model.GetUserInfoByUserName(name)

	if user.Id > 0 {
		ReturnError(c, 4001, "用户名已存在")
		return
	}
	user, err = model.AddUser(name, EncryMd5(pwd))
	if err != nil {
		ReturnError(c, 4001, "添加失败")
		return
	}
	ReturnSuccess(c, 0, "注册成功", user, 1)
}

func (u UserController) Login(c *gin.Context) {
	//获取参数信息
	username := c.DefaultPostForm("name", "")
	password := c.DefaultPostForm("pwd", "")
	if username == "" || password == "" {
		ReturnError(c, 4001, "请输入正确的信息")
		return
	}

	user, _ := model.GetUserInfoByUserName(username)
	if user.Id == 0 || user.Pwd != EncryMd5(password) {
		ReturnError(c, 4001, "用户名或密码不正确")
		return
	}
	data := UserApi{Id: user.Id, Name: user.Name}
	session := sessions.Default(c)
	session.Set("login:"+strconv.Itoa(user.Id), user.Id)
	session.Save()

	ReturnSuccess(c, 0, "success", data, 1)
}
