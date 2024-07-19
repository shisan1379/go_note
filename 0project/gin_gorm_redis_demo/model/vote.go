package model

import (
	"gin_gorm_redis_demo/dao"
	"gorm.io/gorm"
	"time"
)

type Vote struct {
	Id       int   `json:"id"`
	UserId   int   `json:"userId"`
	PlayerId int   `json:"playerId"`
	AddTime  int64 `json:"addTime"`
}

func (Vote) TableName() string {
	return "vote"
}

func AddVote(userId int, playerId int) (int, error) {
	vote := Vote{UserId: userId, PlayerId: playerId, AddTime: time.Now().Unix()}
	err := dao.Db.Create(&vote).Error
	return vote.Id, err
}

func GetVoteInfo(userId int, playerId int) (Vote, error) {
	var vote Vote
	err := dao.Db.Where("user_id = ? AND player_id = ?", userId, playerId).First(&vote).Error
	return vote, err
}

func UpdatePlayerScore(id int) {
	var player Player
	dao.Db.Model(&player).Where("id = ?", id).UpdateColumn("score", gorm.Expr("score + ?", 1))
}
