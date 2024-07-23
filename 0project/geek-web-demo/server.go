package geek_web_demo

import (
	"net"
	"net/http"
)

// 在Go语言中，表达式
// var _ Server = &HttpServer{}
// 是一种类型断言或类型检查的特殊用法，但它实际上并不进行断言操作，而是利用了Go的编译时类型检查特性。
// 这里的 _ 是一个空标识符（blank identifier），它用于忽略变量值，只用于类型检查。
// 确保 HttpServer 实现了 Server 接口
var _ Server = &HttpServer{}

type HandleFunc func(ctx Context)

type Server interface {
	http.Handler //继承 http.Handler 接口
	Start(add string) error

	// AddRoute 增加路由注册的功能
	// method: http 方法
	// path: 路由
	// handler: 业务逻辑
	AddRoute(method string, pattern string, handler HandleFunc)
	//Get(pattern string, handler HandleFunc)
	//Post(pattern string, handler HandleFunc)
}

type HttpServer struct {
	*router
}

//	type HttpSServer struct {
//		HttpServer
//	}

func NewHttpServer() *HttpServer {
	return &HttpServer{router: NewRouter()}
}

// ServeHTTP 处理请求的入口
func (h *HttpServer) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	//框架代码就在这里
	ctx := Context{
		request:  request,
		response: write,
	}
	//接下来匹配路由，并执行业务逻辑
	h.server(ctx)
}

func (h *HttpServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(l, h)
}

func (h *HttpServer) AddRoute(method string, pattern string, handler HandleFunc) {
	//http.HandleFunc(pattern, handler)
}
func (h *HttpServer) Get(pattern string, handler HandleFunc) {
	h.AddRoute(http.MethodGet, pattern, handler)
}

func (h *HttpServer) Post(pattern string, handler HandleFunc) {
	h.AddRoute(http.MethodPost, pattern, handler)
}
func (h *HttpServer) Options(pattern string, handler HandleFunc) {
	h.AddRoute(http.MethodOptions, pattern, handler)
}

func (h *HttpServer) server(ctx Context) {

}
