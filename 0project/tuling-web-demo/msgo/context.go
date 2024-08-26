package msgo

import (
	"encoding/json"
	"errors"
	"fmt"
	validator "github.com/shisan1379/msgo/Validator"
	"github.com/shisan1379/msgo/binding"
	msLog "github.com/shisan1379/msgo/log"
	"github.com/shisan1379/msgo/render"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"unicode"
)

const defaultMaxMemory = 1024 * 32

type Context struct {
	Response              http.ResponseWriter
	Request               *http.Request
	engine                *Engine
	queryCache            url.Values
	formCache             url.Values
	DisallowUnknownFields bool // 参数中必须含有结构体的值
	IsValidate            bool // 是否开启参数校验（按照规则校验参数）
	StatusCode            int
	Logger                *msLog.Logger
}

func checkParam(value reflect.Value, data any, decoder *json.Decoder) error {
	mapData := make(map[string]interface{})
	_ = decoder.Decode(&mapData)
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		required := field.Tag.Get("msgo")
		tag := field.Tag.Get("json")
		value := mapData[tag]
		if value == nil && required == "required" {
			return errors.New(fmt.Sprintf("filed [%s] is required", tag))
		}
	}
	if data != nil {
		marshal, _ := json.Marshal(mapData)
		_ = json.Unmarshal(marshal, data)
	}
	return nil
}
func (c *Context) DealJson(data any) error {
	jsonBinding := binding.JSON
	jsonBinding.DisallowUnknownFields = c.DisallowUnknownFields
	jsonBinding.IsValidate = c.IsValidate
	return c.MustBindWith(data, jsonBinding)
}

func validate(obj any) error {
	return validator.Validator.ValidateStruct(obj)
}

func validateParam(data any, decoder *json.Decoder) error {
	//解析为map，并根据map 中的key 做对比
	//判断类型为 结构体  才能解析为 map

	rVal := reflect.ValueOf(data)

	//是否 为指针
	if rVal.Kind() != reflect.Pointer {
		return errors.New("data is not a pointer")
	}
	elem := rVal.Elem().Interface()

	of := reflect.ValueOf(elem)
	switch of.Kind() {
	case reflect.Struct:
		//将 json 解析为 map
		mapVal := map[string]interface{}{}
		decoder.Decode(&mapVal)
		for i := 0; i < of.NumField(); i++ {
			field := of.Type().Field(i)
			tag := field.Tag.Get("json")
			value := mapVal[tag]
			if value == nil {
				return errors.New(fmt.Sprintf("filed [%s] is not exist", tag))
			}
		}
		//对 map 进行序列化
		marshal, _ := json.Marshal(mapVal)
		// 对 map 进行反序列化，赋值给 data
		_ = json.Unmarshal(marshal, data)

	case reflect.Slice, reflect.Array:
		elem := of.Type().Elem()
		elemType := elem.Kind()
		if elemType == reflect.Struct {
			return checkParamSlice(elem, data, decoder)
		}
	default:
		err := decoder.Decode(data)
		if err != nil {
			return err
		}
	}
	return nil
}
func checkParamSlice(elem reflect.Type, data any, decoder *json.Decoder) error {
	mapData := make([]map[string]interface{}, 0)
	_ = decoder.Decode(&mapData)
	if len(mapData) <= 0 {
		return nil
	}
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Field(i)
		required := field.Tag.Get("msgo")
		tag := field.Tag.Get("json")
		value := mapData[0][tag]
		if value == nil && required == "required" {
			return errors.New(fmt.Sprintf("filed [%s] is required", tag))
		}
	}
	if data != nil {
		marshal, _ := json.Marshal(mapData)
		_ = json.Unmarshal(marshal, data)
	}
	return nil
}
func (c *Context) MultipartForm() (*multipart.Form, error) {
	err := c.Request.ParseMultipartForm(defaultMaxMemory)
	return c.Request.MultipartForm, err
}
func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
func (c *Context) FormFile(name string) (*multipart.FileHeader, error) {
	req := c.Request
	if err := req.ParseMultipartForm(defaultMaxMemory); err != nil {
		return nil, err
	}
	file, header, err := req.FormFile(name)
	if err != nil {
		return nil, err
	}
	err = file.Close()
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		if err := c.Request.ParseMultipartForm(defaultMaxMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				log.Println(err)
			}
		}
		c.formCache = c.Request.PostForm
	}
}

func (c *Context) GetPostForm(key string) (string, bool) {
	if values, ok := c.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) PostFormArray(key string) (values []string) {
	values, _ = c.GetPostFormArray(key)
	return
}

func (c *Context) GetPostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

func (c *Context) GetPostFormMap(key string) (map[string]string, bool) {
	c.initFormCache()
	return c.getMap(c.formCache, key)
}

func (c *Context) PostFormMap(key string) (dicts map[string]string) {
	dicts, _ = c.GetPostFormMap(key)
	return
}
func (c *Context) initQueryCatch() {
	//if c.queryCache == nil {
	//注意此处未做线程安全处理，但是可能不需要
	if c.Request != nil {
		c.queryCache = c.Request.URL.Query()
	} else {
		c.queryCache = url.Values{}
	}
	//}
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
	values = c.queryCache.Get(key)
	return
}
func (c *Context) GetQueryArray(key string) (values []string, ok bool) {
	c.initQueryCatch()
	values, ok = c.queryCache[key]
	return
}

func (c Context) GetQueryMap(key string) (map[string]string, bool) {
	c.initQueryCatch()
	return c.getMap(c.queryCache, key)
}

func (c Context) getMap(m map[string][]string, key string) (map[string]string, bool) {
	//user[id]=1&user[name]=张三
	dicts := make(map[string]string)
	exist := false
	for k, val := range m {
		i := strings.IndexByte(k, '[')
		if i >= 1 && k[0:i] == key {
			j := strings.IndexByte(k[i+1:], ']')
			if j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = val[0]
			}

		}
	}
	return dicts, exist
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
func (c *Context) BindXML(obj any) error {
	return c.MustBindWith(obj, binding.XML)
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
	c.StatusCode = code
	err := r.Render(c.Response)
	return err
}

func (c *Context) MustBindWith(data any, b binding.Binding) error {
	//如果发生错误，返回400状态码 参数错误
	if err := c.ShouldBindWith(data, b); err != nil {
		c.Response.WriteHeader(http.StatusBadRequest)
		return err
	}
	return nil
}
func (c *Context) ShouldBindWith(obj any, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}

func (c *Context) Fail(code int, s string) {
	c.String(code, s)
}

func (c *Context) HandlerWithError(err error) {
	if err != nil {
		code, data := c.engine.ErrHandler(err)
		c.Json(code, data)
	}

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
