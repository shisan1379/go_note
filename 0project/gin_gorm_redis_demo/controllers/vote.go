package controllers

import (
	cache "gin_gorm_redis_demo/catch"
	"gin_gorm_redis_demo/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type VoteController struct{}

func (r VoteController) AddVote(c *gin.Context) {
	userIdStr := c.DefaultPostForm("userId", "0")
	playerIdStr := c.DefaultPostForm("playerId", "0")
	userId, _ := strconv.Atoi(userIdStr)
	playerId, _ := strconv.Atoi(playerIdStr)

	if userId == 0 || playerId == 0 {
		ReturnError(c, 4001, "请输入正确的信息")
		return
	}
	user, _ := model.GetUserInfo(userId)
	if user.Id == 0 {
		ReturnError(c, 4001, "投票用户不存在")
		return
	}
	player, _ := model.GetPlayerInfo(playerId)
	if player.Id == 0 {
		ReturnError(c, 4001, "选手不存在")
		return
	}
	vote, _ := model.GetVoteInfo(userId, playerId)
	if vote.Id != 0 {
		ReturnError(c, 4001, "已投票")
		return
	}

	rs, err := model.AddVote(userId, playerId)
	if err == nil {
		//更新选手票数
		model.UpdatePlayerScore(playerId)
		//更新redis
		var redisKey string
		redisKey = "ranking:" + strconv.Itoa(player.Aid)
		cache.Rdb.ZIncrBy(cache.Rctx, redisKey, 1, strconv.Itoa(playerId))
		ReturnSuccess(c, 0, "success", rs, 1)
		return
	}
	ReturnError(c, 4004, "没有相关信息")

}
