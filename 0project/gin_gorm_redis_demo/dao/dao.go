package dao

import (
	"gin_gorm_redis_demo/config"
	logger "gin_gorm_redis_demo/pkg/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	Db  *gorm.DB
	err error
)

func init() {
	newLogger := glog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		glog.Config{
			SlowThreshold:             time.Second, // 慢速 SQL 阈值
			LogLevel:                  glog.Info,   // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略记录器的 ErrRecordNotFound 错误
			ParameterizedQueries:      false,       // 不要在 SQL 日志中包含参数
			Colorful:                  false,       // 禁用颜色
		},
	)
	gConfig := gorm.Config{
		Logger: newLogger,
	}

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
