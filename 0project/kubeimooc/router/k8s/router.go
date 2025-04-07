package k8s

import (
	"github.com/gin-gonic/gin"
	"kubeimooc/api"
)

type InitK8sRouter struct {
}

func (router *InitK8sRouter) InitK8sRouter(r *gin.Engine) {
	group := r.Group("/k8s")
	apiGroup := api.ApiGroupApp.K8sApiGroup

	group.GET("/listPod", apiGroup.GetPodList)
	group.GET("/namespace", apiGroup.GetNamespaceList)
}
