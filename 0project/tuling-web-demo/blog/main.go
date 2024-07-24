package main

import (
	"fmt"
	"github.com/shisan1379/msgo"
)

func main() {

	engine := msgo.New()
	group := engine.Group("user")
	group.PreHandle(func(next msgo.HandleFunc) msgo.HandleFunc {
		return func(ctx *msgo.Context) {
			fmt.Println("pre handler")
			next(ctx)
		}
	})
	group.Get("/hello", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s 欢迎来到我的世界", "get")
	})

	group.Post("/hello", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s 欢迎来到我的世界", "post")
	})

	group.Get("/get/:id", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s get", ":id")
	})

	group.Get("/hello/*/get", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s ", "/hello/*get")
	})

	engine.Run()
}
