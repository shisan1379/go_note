package geek_web_demo

import (
	"fmt"
	"testing"
)

func Test_router_AddRoute(t *testing.T) {
	//1. 构造路由树
	//2. 验证路由树
	testRoutes := []struct {
		method string
		path   string
	}{}
	var mockHandler HandleFunc = func(ctx Context) {

	}
	r := NewRouter()
	for _, route := range testRoutes {
		r.AddRoute(route.method, route.path, mockHandler)
	}

	wantRouter := &router{
		Trees: map[string]*node{},
	}
	msg, ok := wantRouter.equal(*r)
	if !ok {
	}

}
func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.Trees {
		dst, ok := y.Trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http method"), false
		}
		msg, r := v.equal(dst)
		if !r {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {

}
