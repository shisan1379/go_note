package msgo

import (
	"fmt"
	"github.com/shisan1379/msgo/render"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"unicode"
)

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	engine     *Engine
	queryCatch url.Values
}

func (c *Context) initQueryCatch() {
	if c.queryCatch == nil {
		//注意此处未做线程安全处理，但是可能不需要
		if c.Request != nil {
			c.queryCatch = c.Request.URL.Query()
		} else {
			c.queryCatch = url.Values{}
		}
	}
}

func (c *Context) GetDefaultQuery(key string, defaultVal string) (values string) {
	array, ok := c.GetQueryArray(key)
	if ok {
		return array[0]
	}
	return defaultVal
}

func (c *Context) GetDefaultQueryArray(key string, defaultVal []string) (values []string) {
	array, ok := c.GetQueryArray(key)
	if ok {
		return array
	}
	return defaultVal
}

func (c *Context) GetQuery(key string) (values string) {
	c.initQueryCatch()
	values = c.queryCatch.Get(key)
	return
}
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCatch()
	values, ok = c.queryCatch[key]
	return
}

func (ctx *Context) HTML(status int, html string) error {
	return ctx.Render(status, render.HTML{IsTemplate: false, Data: html})
}

func (c *Context) HTMLTemplate(name string, data any, fileName ...string) {
	c.Render(http.StatusOK, render.HTML{
		IsTemplate: true,
		Name:       name,
		Data:       data,
		Template:   c.engine.HTMLRender.Template,
	})
}
func (c *Context) HTMLTemplateGlob2(name string, data any, pattern string) {
	t := template.New(name)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		log.Println(err)
		return
	}
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = t.Execute(c.Response, data)

	//err = t.ExecuteTemplate(c.Response, name, data)
	// 执行模板渲染，将结果写入标准输出（或其他 io.Writer）
	//err = t.ExecuteTemplate(os.Stdout, "login.html", data)
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

func (c *Context) Template(name string, data any) error {
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.engine.HTMLRender.Template.ExecuteTemplate(c.Response, name, data)
	return err
}

func (c *Context) Json(status int, data any) error {
	err := c.Render(status, &render.JSON{Data: data})
	return err

}
func (c *Context) Xml(status int, data any) error {
	return c.Render(status, render.XML{Data: data})

}
func (c *Context) File(name string) {
	http.ServeFile(c.Response, c.Request, name)
}
func (c *Context) FileAttachment(path string, newName string) {
	if isASCII(newName) {
		c.Response.Header().Set("Content-Disposition", `attachment; filename="`+newName+`"`)
	} else {
		c.Response.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(newName))
	}
	http.ServeFile(c.Response, c.Request, path)
}

// path 是相对文件系统的路径
func (c *Context) FileFromFS(path string, fs http.FileSystem) {
	defer func(old string) {
		fmt.Println("路径：", old)
		c.Request.URL.Path = old
	}(c.Request.URL.Path)
	fmt.Println(c.Request.URL.Path)
	c.Request.URL.Path = path

	http.FileServer(fs).ServeHTTP(c.Response, c.Request)
}

func (c *Context) Redirect(status int, url string) {
	http.Redirect(c.Response, c.Request, url, status)

	//c.Render(status, render.Redirect{
	//	Code:     status,
	//	Request:  c.Request,
	//	Location: url,
	//})

}

func (c *Context) String(status int, format string, values ...any) error {
	err := c.Render(status, render.String{
		Format: format,
		Data:   values,
	})
	return err
}

func (c *Context) Render(code int, r render.Render) error {
	c.Response.WriteHeader(code)
	err := r.Render(c.Response)
	return err
}

// 是否 ASCII 字符
func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}
