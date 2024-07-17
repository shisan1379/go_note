package controllers

import "github.com/gin-gonic/gin"

type Order struct {
}

// 通过 key 从 post 中获取参数
func (o *Order) GetList(c *gin.Context) {
	cid := c.PostForm("cid")
	name := c.DefaultPostForm("name", "王五")
	ReturnSuccess(c, 0, cid, name, 10)
}

// 使用 map 接收参数
func (o *Order) GetList2(c *gin.Context) {

	param := make(map[string]interface{})
	err := c.BindJSON(&param)
	if err != nil {
		ReturnError(c, 1, err.Error())
		return
	}
	ReturnSuccess(c, 0, param["cid"], param["name"], 10)
}

type Search struct {
	Cid  int    `form:"cid"`
	Name string `form:"name"`
}

// 使用 struct 接收参数
func (o *Order) SearchList(c *gin.Context) {

	//处理可能发生的异常
	defer func() {
		if err := recover(); err != nil {
			ReturnError(c, 1, "异常")
		}
	}()

	search := Search{}

	err := c.BindJSON(&search)
	if err != nil {
		ReturnError(c, 1, err.Error())
		return
	}
	ReturnSuccess(c, 0, search.Cid, search.Name, 10)
}
