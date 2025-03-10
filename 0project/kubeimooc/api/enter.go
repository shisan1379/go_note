package api

import "kubeimooc/api/example"

type ApiGroup struct {
	ExampleApiGroup example.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
