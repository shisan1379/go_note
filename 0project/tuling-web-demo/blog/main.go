package main

import (
	"fmt"
	"github.com/shisan1379/msgo"
	"net/http"
)

type User struct {
	Name string
}

func Log(next msgo.HandleFunc) msgo.HandleFunc {
	return func(ctx *msgo.Context) {
		fmt.Println("pre log")
		next(ctx)
		fmt.Println("post log")
	}
}
func main() {

	engine := msgo.New()
	group := engine.Group("user")

	//前置中间件
	group.Use(func(next msgo.HandleFunc) msgo.HandleFunc {
		return func(ctx *msgo.Context) {
			fmt.Println("pre handler")
			next(ctx)
			fmt.Println("post handler")
		}
	})

	//group.PostHandle(func(next msgo.HandleFunc) msgo.HandleFunc {
	//	return func(ctx *msgo.Context) {
	//		fmt.Println("post handler")
	//	}
	//})

	group.Get("/hello", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s 欢迎来到我的世界", "get")
	}, Log)

	group.Post("/hello", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s 欢迎来到我的世界", "post")
	}, Log)

	group.Get("/get/:id", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s get", ":id")
	}, Log)

	group.Get("/hello/*/get", func(ctx *msgo.Context) {
		//写入到标准输出 -> w
		// w -> 前端
		fmt.Fprintf(ctx.Response, "%s ", "/hello/*get")
	}, Log)
	group.Get("/html", func(ctx *msgo.Context) {
		ctx.HTML(http.StatusOK, "<h1>啦啦啦啦</h1>")
	})

	group.Get("/template", func(ctx *msgo.Context) {
		ctx.HTMLTemplate("index", "", "tpl/index.html")
	})
	group.Get("/login", func(ctx *msgo.Context) {
		ctx.HTMLTemplate("login", &User{Name: "123"}, "tpl/login.html", "tpl/header.html")
	})
	group.Get("/login2", func(ctx *msgo.Context) {
		ctx.HTMLTemplateGlob2("login.html", &User{Name: "123123123"}, "tpl/*.html")
	})

	engine.Run()
}
