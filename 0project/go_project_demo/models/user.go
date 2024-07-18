package models

import "go_project_demo/dao/dao"

type User struct {
	Id   int
	Name string
}

func (User) TableName() string {
	return "user"
}

func GetUserTest(id int) (User, error) {
	var user User
	err := dao.Db.Where("id = ?", id).First(&user).Error
	return user, err
}

func AddUser(name string) (User, error) {
	user := User{Name: name}
	err := dao.Db.Create(&user).Error
	return user, err
}
