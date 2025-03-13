package api

import (
	"kubeimooc/api/example"
	"kubeimooc/api/k8s"
)

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
	K8sApiGroup     k8s.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
