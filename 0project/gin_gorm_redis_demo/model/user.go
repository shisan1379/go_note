package model

import "gin_gorm_redis_demo/dao"

type User struct {
	Id       int `gorm:"primary_key;AUTO_INCREMENT"`
	Username string
	Password string
}

func (User) TableName() string {
	return "user"
}

func GetUserInfoByUserName(userName string) (User, error) {
	var user User
	err := dao.Db.Where("name = ?", userName).First(&user).Error
	return user, err
}

func AddUser(name string, pwd string) (user User, err error) {
	user.Username = name
	user.Password = pwd
	err = dao.Db.Create(&user).Error
	return
}

func GetUserInfo(id int) (User, error) {
	var user User
	err := dao.Db.Where("id = ?", id).First(&user).Error
	return user, err
}
