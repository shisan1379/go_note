package router

import "kubeimooc/router/example"

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
}

var RouterGroupApp = new(RouterGroup)
