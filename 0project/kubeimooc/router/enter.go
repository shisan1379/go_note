package router

import (
	"kubeimooc/router/example"
	"kubeimooc/router/k8s"
)

type RouterGroup struct {
	ExampleRouterGroup example.ExampleRouter
	K8sRouterGroup     k8s.InitK8sRouter
}

var RouterGroupApp = new(RouterGroup)
