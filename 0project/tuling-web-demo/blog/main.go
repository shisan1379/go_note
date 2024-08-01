package main

import (
	"fmt"
	"github.com/shisan1379/msgo"
	"net/http"
)

type User struct {
	Name string `xml:"name"`
	Age  int    `xml:"age"`
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

	//group.Get("/template", func(ctx *msgo.Context) {
	//	ctx.HTMLTemplate("index", "", "tpl/index.html")
	//})
	//group.Get("/login", func(ctx *msgo.Context) {
	//	ctx.HTMLTemplate("login", &User{Name: "123"}, "tpl/login.html", "tpl/header.html")
	//})
	//group.Get("/login2", func(ctx *msgo.Context) {
	//	ctx.HTMLTemplateGlob2("login.html", &User{Name: "123123123"}, "tpl/*.html")
	//})
	//group.Get("/test", func(ctx *msgo.Context) {
	//	ctx.HTMLTemplateGlob2("test.html", nil, "tpl/*.html")
	//})
	engine.LoadTemplate("tpl/*.html")
	group.Get("/login", func(ctx *msgo.Context) {
		ctx.Template("login.html", &User{Name: "123123123"})
	})
	group.Get("/json", func(ctx *msgo.Context) {
		ctx.Json(http.StatusOK, &User{Name: "123123123"})
	})
	group.Get("/xml", func(ctx *msgo.Context) {
		ctx.Xml(http.StatusOK, &User{Name: "123123123", Age: 10})
	})
	group.Get("/file", func(ctx *msgo.Context) {
		ctx.File("./tpl/bb.xlsx")
	})
	group.Get("/fileName", func(ctx *msgo.Context) {
		ctx.FileAttachment("./tpl/bb.xlsx", "aaa.xlsx")
	})
	group.Get("/fs", func(ctx *msgo.Context) {
		ctx.FileFromFS("bb.xlsx", http.Dir("tpl"))
	})
	group.Get("/redirect", func(ctx *msgo.Context) {
		ctx.Redirect(http.StatusFound, "user/login")
	})
	group.Get("/string", func(ctx *msgo.Context) {
		ctx.String(http.StatusOK, "string")
	})
	group.Get("/add", func(ctx *msgo.Context) {
		ids := ctx.GetQuery("id")
		fmt.Println(ids)
		ctx.String(http.StatusOK, "string")
	})
	group.Get("/add2", func(ctx *msgo.Context) {
		ids, _ := ctx.GetQueryArray("ids")
		fmt.Println(ids)
		ctx.String(http.StatusOK, "string")
	})
	engine.Run()
}
