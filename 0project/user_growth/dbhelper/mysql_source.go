package dbhelper

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"user_growth/conf"
	"xorm.io/xorm"
)

var dbEngine *xorm.Engine

func InitDb() {
	if dbEngine != nil {
		return
	}

	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		conf.GlobalConfig.DB.Username,
		conf.GlobalConfig.DB.Password,
		conf.GlobalConfig.DB.Host,
		conf.GlobalConfig.DB.Port,
		conf.GlobalConfig.DB.Database,
		conf.GlobalConfig.DB.Charset)

	engine, err := xorm.NewEngine(conf.GlobalConfig.DB.Engine, sourceName)
	if err != nil {
		log.Fatalf("dbhelper NewEngine err:%v", err)
	} else {
		dbEngine = engine
	}
	if conf.GlobalConfig.DB.MaxIdleConns > 0 {
		dbEngine.SetMaxIdleConns(conf.GlobalConfig.DB.MaxIdleConns)
	}

	if conf.GlobalConfig.DB.MaxOpenConns > 0 {
		dbEngine.SetMaxOpenConns(conf.GlobalConfig.DB.MaxOpenConns)
	}

	if conf.GlobalConfig.DB.ConnMaxLifetime > 0 {
		dbEngine.SetConnMaxLifetime(time.Minute * time.Duration(conf.GlobalConfig.DB.ConnMaxLifetime))
	}
}

func GetDb() *xorm.Engine {
	return dbEngine
}
