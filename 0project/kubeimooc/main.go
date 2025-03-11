package main

import (
	"kubeimooc/global"
	"kubeimooc/initiallize"
)

// 项目的启动入口
func main() {
	initiallize.Viper()
	r := initiallize.Routers()
	initiallize.K8S()
	panic(r.Run(global.CONF.System.Addr))

}
