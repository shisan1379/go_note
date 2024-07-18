package dao

import (
	"go_project_demo/config"
	logger "go_project_demo/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	Db  *gorm.DB
	err error
)

func init() {

	gConfig := gorm.Config{}

	Db, err = gorm.Open(mysql.Open(config.Mysqldb), &gConfig)
	if err != nil {

		logger.Error(map[string]interface{}{"mysql connect error": err.Error()})
	}
	if Db.Error == nil {
		logger.Error(map[string]interface{}{"database connect error": Db.Error})
	}

	// 获取通用数据库对象 sql.DB 以便设置连接池参数
	sqlDB, _ := Db.DB()

	// 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	//  设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
