package dbhelper

import (
	"fmt"
	"log"
	"user_growth/conf"
	"xorm.io/xorm"
)

var dbEngine *xorm.Engine

func InitDb() {
	if dbEngine != nil {
		return
	}

	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		conf.GlobalConfig.DB.UserName,
		conf.GlobalConfig.DB.Password,
		conf.GlobalConfig.DB.Host,
		conf.GlobalConfig.DB.Port,
		conf.GlobalConfig.DB.Database,
		conf.GlobalConfig.DB.Charset)

	engine, err := xorm.NewEngine(conf.GlobalConfig.DB.Engine, sourceName)
	if err != nil {
		log.Fatalf("dbhelper NewEngine err:%v", err)
	}

}
