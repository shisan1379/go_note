package msgo

import (
	"fmt"
	"log"
	"net/http"
)

const ANY = "ANY"

type HandleFunc func(ctx *Context)

type MiddleWare func(next HandleFunc) HandleFunc

type routerGroup struct {
	name              string
	handleFuncMap     map[string]map[string]HandleFunc   //<path,<method, HandleFunc>>
	middleWareFuncMap map[string]map[string][]MiddleWare //<path,<method, []MiddleWare>>
	routerTree        *treeNode
	preMiddleWare     []MiddleWare
	postMiddleWare    []MiddleWare
}

// Use 添加中间件
func (r *routerGroup) Use(middleWareFunc ...MiddleWare) {
	r.preMiddleWare = append(r.preMiddleWare, middleWareFunc...)
}

// PostHandle 添加后置中间件
//func (r *routerGroup) PostHandle(middleWareFunc ...MiddleWare) {
//	r.postMiddleWare = append(r.postMiddleWare, middleWareFunc...)
//}

// MethodHandle 执行方法处理，包含执行中间件
func (r *routerGroup) MethodHandle(routerPath string, method string, handleFunc HandleFunc, ctx *Context) {
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
}

// Group 为路由添加组
func (receiver *router) Group(name string) *routerGroup {
	// 这里没有检查组是否重复
	group := routerGroup{
		name:              name,
		handleFuncMap:     make(map[string]map[string]HandleFunc),
		middleWareFuncMap: make(map[string]map[string][]MiddleWare),
		routerTree: &treeNode{
			name:     "/",
			children: make([]*treeNode, 0),
		},
	}
	receiver.routerGroups = append(receiver.routerGroups, &group)
	return &group
}

// addRouter 添加路由
func (r routerGroup) addRouter(method string, path string, handleFunc HandleFunc, ware ...MiddleWare) {

	m, _ := r.handleFuncMap[path]
	if m == nil {
		r.handleFuncMap[path] = make(map[string]HandleFunc)
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

func (r routerGroup) Any(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(ANY, name, handleFunc, ware...)
}

func (r routerGroup) Get(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodGet, name, handleFunc, ware...)
}
func (r routerGroup) Post(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodPost, name, handleFunc, ware...)
}
func (r routerGroup) Delete(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodDelete, name, handleFunc, ware...)
}
func (r routerGroup) Put(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodPut, name, handleFunc, ware...)
}

func (r routerGroup) Patch(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodPatch, name, handleFunc, ware...)
}

func (r routerGroup) Options(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodOptions, name, handleFunc, ware...)
}

func (r routerGroup) Head(name string, handleFunc HandleFunc, ware ...MiddleWare) {
	r.addRouter(http.MethodHead, name, handleFunc, ware...)
}

// Engine  引擎
type Engine struct {
	router
}

// ServeHTTP 实现处理方法
func (e *Engine) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	e.httpRequestHandle(write, request)
}

// httpRequestHandle http请求处理器
func (e *Engine) httpRequestHandle(write http.ResponseWriter, request *http.Request) {
	method := request.Method
	for _, group := range e.routerGroups {

		routerName := SubStringLast(request.RequestURI, "/"+group.name)

		node := group.routerTree.Get(routerName)
		if node != nil && node.isEnd {

			//路由匹配上了
			ctx := &Context{
				Response: write,
				Request:  request,
			}
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.MethodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.MethodHandle(node.routerName, method, handle, ctx)
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
	return &Engine{
		router: router{
			routerGroups: []*routerGroup{},
		},
	}
}

func (e *Engine) Run() {
	//for _, g := range e.routerGroups {
	//	for key, val := range g.handleFuncMap {
	//		http.HandleFunc("/"+g.name+key, val)
	//	}
	//}
	http.Handle("/", e)

	err := http.ListenAndServe(":8111", nil)

	if err != nil {
		log.Fatal(err)
	}
}
