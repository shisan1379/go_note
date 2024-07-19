package model

import "gin_gorm_redis_demo/dao"

type Player struct {
	Id          int    `json:"id"`
	Aid         int    `json:"aid"`
	Ref         string `json:"ref"`
	NickName    string `json:"nick_name"`
	Declaration string `json:"declaration"`
	Avatar      string `json:"avatar"`
	Score       int    `json:"score"`
}

func (receiver Player) TableName() string {
	return "player"
}

func GetPlayers(aid int) ([]Player, error) {
	var player []Player
	err := dao.Db.Where("aid = ?", aid).Find(&player).Error

	return player, err
}

func GetPlayerInfo(id int) (Player, error) {
	var player Player
	err := dao.Db.Where("id = ?", id).First(&player).Error
	return player, err
}
