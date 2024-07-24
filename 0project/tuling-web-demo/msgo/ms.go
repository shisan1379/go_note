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
	name           string
	handleFuncMap  map[string]map[string]HandleFunc //<path,<method, HandleFunc>>
	routerTree     *treeNode
	preMiddleWare  []MiddleWare
	postMiddleWare []MiddleWare
}

func (r *routerGroup) PreHandle(middleWareFunc ...MiddleWare) {
	r.preMiddleWare = append(r.preMiddleWare, middleWareFunc...)
}
func (r *routerGroup) PostHandle(middleWareFunc ...MiddleWare) {
	r.postMiddleWare = append(r.postMiddleWare, middleWareFunc...)
}
func (r *routerGroup) MethodHandle(handleFunc HandleFunc, ctx *Context) {
	if r.preMiddleWare != nil {
		for _, ware := range r.preMiddleWare {
			handleFunc = ware(handleFunc)
		}
	}
	handleFunc(ctx)
	if r.postMiddleWare != nil {
		for _, ware := range r.postMiddleWare {
			handleFunc = ware(handleFunc)
		}
	}
}

type router struct {
	routerGroups []*routerGroup
}

func (receiver *router) Group(name string) *routerGroup {
	group := routerGroup{
		name:          name,
		handleFuncMap: make(map[string]map[string]HandleFunc),
		routerTree: &treeNode{
			name:     "/",
			children: make([]*treeNode, 0),
		},
	}
	receiver.routerGroups = append(receiver.routerGroups, &group)
	return &group
}

func (r routerGroup) addRouter(method string, path string, handleFunc HandleFunc) {

	m, _ := r.handleFuncMap[path]
	if m == nil {
		r.handleFuncMap[path] = make(map[string]HandleFunc)
	}
	_, ok := r.handleFuncMap[path][method]
	if ok {
		panic(fmt.Sprintf("%s /%s%s 同一个路由下不能重复", method, r.name, path))
	}

	r.handleFuncMap[path][method] = handleFunc
	r.routerTree.Put(path)

}
func (r routerGroup) Any(name string, handleFunc HandleFunc) {
	r.addRouter(ANY, name, handleFunc)
}

func (r routerGroup) Get(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodGet, name, handleFunc)
}
func (r routerGroup) Post(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodPost, name, handleFunc)
}
func (r routerGroup) Delete(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodDelete, name, handleFunc)
}
func (r routerGroup) Put(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodPut, name, handleFunc)
}

func (r routerGroup) Patch(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodPatch, name, handleFunc)
}

func (r routerGroup) Options(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodOptions, name, handleFunc)
}

func (r routerGroup) Head(name string, handleFunc HandleFunc) {
	r.addRouter(http.MethodHead, name, handleFunc)
}

// Engine  引擎
type Engine struct {
	router
}

func (e *Engine) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	e.httpRequestHandle(write, request)
}

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
				group.MethodHandle(handle, ctx)
				return
			}
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.MethodHandle(handle, ctx)
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
