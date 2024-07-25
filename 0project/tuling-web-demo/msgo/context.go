package msgo

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
}

func (ctx *Context) HTML(status int, html string) error {
	// 状态码默认为 200
	ctx.Response.WriteHeader(status)
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := ctx.Response.Write([]byte(html))
	return err
}

func (c *Context) HTMLTemplate(name string, data any, fileName ...string) {
	t := template.New(name)
	t, err := t.ParseFiles(fileName...)
	if err != nil {
		log.Println(err)
		return
	}
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.Response, data)
	if err != nil {
		log.Println(err)
	}
}
func (c *Context) HTMLTemplateGlob2(name string, data any, pattern string) {
	t := template.New(name)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	//err = t.Execute(c.Response, data)
	// 执行模板渲染，将结果写入标准输出（或其他 io.Writer）
	err = t.ExecuteTemplate(os.Stdout, "login.html", data)
	if err != nil {
		log.Println(err)
	}
}

func (c *Context) HTMLTemplateGlob(name string, funcMap template.FuncMap, pattern string, data any) {
	t := template.New(name)
	t.Funcs(funcMap)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.Response, data)
	if err != nil {
		log.Println(err)
	}
}
