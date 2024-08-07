package main

import (
	"fmt"
	"github.com/shisan1379/msgo"
	"io"
	"log"
	"net/http"
	"os"
)

type User struct {
	Name      string   `xml:"name" json:"name" msgo:"required"`
	Age       int      `xml:"name" json:"age" validate:"required,max=50,min=18"`
	Addresses []string `json:"addresses"`
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
	group.Get("/add3", func(ctx *msgo.Context) {
		name := ctx.GetDefaultQuery("name", "张三")
		fmt.Println(name)
		ctx.String(http.StatusOK, name)
	})
	group.Get("/map", func(ctx *msgo.Context) {
		name, _ := ctx.GetQueryMap("user")
		fmt.Println(name)

		ctx.Json(http.StatusOK, name)
	})
	group.Post("/form", func(ctx *msgo.Context) {
		name, _ := ctx.GetPostForm("name")
		ctx.Json(http.StatusOK, name)
	})

	group.Post("/upFile", func(ctx *msgo.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			log.Println(err)
		}
		src, err := file.Open()
		defer src.Close()
		if err != nil {
			log.Println(err)
		} else {
			out, err := os.Create("d:/aa/" + file.Filename)
			defer out.Close()
			if err != nil {
				log.Println(err)
			} else {
				io.Copy(out, src)
			}
		}
	})

	group.Post("/jsonParam", func(ctx *msgo.Context) {
		user := &User{}
		err := ctx.DealJson(user)
		if err == nil {
			ctx.Json(http.StatusOK, user)
		} else {
			log.Println(err)
		}
	})

	group.Post("/jsonArray", func(ctx *msgo.Context) {
		user := make([]User, 10)
		err := ctx.DealJson(&user)
		if err == nil {
			ctx.Json(http.StatusOK, user)
		} else {
			log.Println(err)
		}
	})

	engine.Run()
}
