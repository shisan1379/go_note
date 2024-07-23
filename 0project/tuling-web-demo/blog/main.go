package main

import (
	"fmt"
	"github.com/shisan1379/msgo"
	"net/http"
)

func main() {
	//http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
	//	fmt.Fprintf(writer, "%s 欢迎来到我的世界", "you")
	//})
	//err := http.ListenAndServe(":8111", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	engine := msgo.New()
	group := engine.Group("user")
	group.AddRouter("/hello", func(w http.ResponseWriter, r *http.Request) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(w, "%s 欢迎来到我的世界", "you")

	})
	engine.Run()
}
