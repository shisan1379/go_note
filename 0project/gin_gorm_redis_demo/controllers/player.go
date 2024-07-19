package controllers

import (
	"gin_gorm_redis_demo/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PlayerController struct{}

func (p PlayerController) GetPlayers(c *gin.Context) {
	aidStr := c.DefaultPostForm("aid", "0")
	aid, _ := strconv.Atoi(aidStr)

	players, err := model.GetPlayers(aid)
	if err != nil {
		ReturnError(c, 4004, err.Error())
		return
	}

	ReturnSuccess(c, 0, "success", players, int64(len(players)))
}
