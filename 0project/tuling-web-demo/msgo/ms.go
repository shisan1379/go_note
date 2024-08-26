package msgo

import (
	"fmt"
	msLog "github.com/shisan1379/msgo/log"
	"github.com/shisan1379/msgo/render"
	"html/template"
	"log"
	"net/http"
	"sync"
)

const ANY = "ANY"

type HandlerFunc func(ctx *Context)

type MiddleWare func(next HandlerFunc) HandlerFunc

type routerGroup struct {
	name              string
	handleFuncMap     map[string]map[string]HandlerFunc  //<path,<method, HandlerFunc>>
	middleWareFuncMap map[string]map[string][]MiddleWare //<path,<method, []MiddleWare>>
	routerTree        *treeNode
	preMiddleWare     []MiddleWare
}

// Use 添加中间件
func (r *routerGroup) Use(middleWareFunc ...MiddleWare) {
	r.preMiddleWare = append(r.preMiddleWare, middleWareFunc...)
}

// MethodHandle 执行方法处理，包含执行中间件
func (r *routerGroup) MethodHandle(routerPath string, method string, handleFunc HandlerFunc, ctx *Context) {
	// 路由级别的中间件
	middleWareFuncs, ok := r.middleWareFuncMap[routerPath][method]
	if ok {
		for _, wareFunc := range middleWareFuncs {
			handleFunc = wareFunc(handleFunc)
		}
	}

	//通用级别的中间件
	if r.preMiddleWare != nil {
		for _, ware := range r.preMiddleWare {
			handleFunc = ware(handleFunc)
		}
	}
	handleFunc(ctx)
}

// router 路由
type router struct {
	routerGroups []*routerGroup
	engine       *Engine
}

// Group 为路由添加组
func (receiver *router) Group(name string) *routerGroup {
	// 这里没有检查组是否重复
	group := routerGroup{
		name:              name,
		handleFuncMap:     make(map[string]map[string]HandlerFunc),
		middleWareFuncMap: make(map[string]map[string][]MiddleWare),
		routerTree: &treeNode{
			name:     "/",
			children: make([]*treeNode, 0),
		},
	}
	group.Use(receiver.engine.Middles...)
	receiver.routerGroups = append(receiver.routerGroups, &group)
	return &group
}

// addRouter 添加路由
func (r routerGroup) addRouter(method string, path string, handleFunc HandlerFunc, ware ...MiddleWare) {

	m, _ := r.handleFuncMap[path]
	if m == nil {
		r.handleFuncMap[path] = make(map[string]HandlerFunc)
		r.middleWareFuncMap[path] = make(map[string][]MiddleWare)
	}
	_, ok := r.handleFuncMap[path][method]
	if ok {
		panic(fmt.Sprintf("%s /%s%s 同一个路由下不能重复", method, r.name, path))
	}
	if ware != nil && len(ware) > 0 {
		r.middleWareFuncMap[path][method] = append(r.middleWareFuncMap[path][method], ware...)
	}

	r.handleFuncMap[path][method] = handleFunc
	r.routerTree.Put(path)

}

func (r routerGroup) Any(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(ANY, name, handleFunc, ware...)
}

func (r routerGroup) Get(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodGet, name, handleFunc, ware...)
}
func (r routerGroup) Post(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodPost, name, handleFunc, ware...)
}
func (r routerGroup) Delete(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodDelete, name, handleFunc, ware...)
}
func (r routerGroup) Put(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodPut, name, handleFunc, ware...)
}

func (r routerGroup) Patch(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodPatch, name, handleFunc, ware...)
}

func (r routerGroup) Options(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodOptions, name, handleFunc, ware...)
}

func (r routerGroup) Head(name string, handleFunc HandlerFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodHead, name, handleFunc, ware...)
}

type ErrorHandler func(err error) (code int, msg any)

// Engine  引擎
type Engine struct {
	router
	funcMap    template.FuncMap
	HTMLRender render.HTMLRender
	pool       sync.Pool
	Logger     *msLog.Logger
	Middles    []MiddleWare
	ErrHandler ErrorHandler
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}
func (e *Engine) LoadTemplate(pattern string) {
	tmpl := template.Must(
		template.New("").
			Funcs(e.funcMap).
			ParseGlob(pattern))
	e.SetTemplate(tmpl)
}

func (e *Engine) RegisterErrorhandler(eh ErrorHandler) {
	e.ErrHandler = eh
}

func (e *Engine) SetTemplate(t *template.Template) {
	e.HTMLRender = render.HTMLRender{Template: t}
}

// ServeHTTP 实现处理方法
func (e *Engine) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	context := e.pool.Get().(*Context)
	context.Response = write
	context.Request = request
	context.Logger = e.Logger

	e.httpRequestHandle(context, write, request)

	e.pool.Put(context)
}

// httpRequestHandle http请求处理器
func (e *Engine) httpRequestHandle(context *Context, write http.ResponseWriter, request *http.Request) {
	method := request.Method
	for _, group := range e.routerGroups {

		routerName := SubStringLast(request.URL.Path, "/"+group.name)

		node := group.routerTree.Get(routerName)
		if node != nil && node.isEnd {

			//路由匹配上了
			//ctx := &Context{
			//	Response: write,
			//	Request:  request,
			//	engine:   e,
			//}

			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.MethodHandle(node.routerName, ANY, handle, context)
				return
			}
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.MethodHandle(node.routerName, method, handle, context)
				return
			}
		}
		write.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(write, "%s %s not allowed \n", request.RequestURI, method)
		return
	}
	fmt.Fprintf(write, "404")
}

func New() *Engine {
	engine := &Engine{
		router: router{
			routerGroups: []*routerGroup{},
		},
	}
	engine.pool.New = func() any {
		return engine.allocateContext()
	}
	return engine

}

func Default() *Engine {
	engine := New()
	engine.Logger = msLog.Default()
	engine.Use(Logging)  //默认使用log 中间件
	engine.Use(Recovery) //默认使用 recover中间件
	engine.router.engine = engine
	return engine

}

func (e *Engine) allocateContext() any {
	return &Context{engine: e}
}

func (e *Engine) Run() {
	//for _, g := range e.routerGroups {
	//	for key, val := range g.handleFuncMap {
	//		http.HandlerFunc("/"+g.name+key, val)
	//	}
	//}
	http.Handle("/", e)

	err := http.ListenAndServe(":8111", nil)

	if err != nil {
		log.Fatal(err)
	}
}

func (e *Engine) Use(middles ...MiddleWare) {
	e.Middles = middles
}
