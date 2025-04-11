package k8s

import (
	"context"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeimooc/global"
	namespace_res "kubeimooc/model/namespace/response"
	"kubeimooc/response"
)

type NamespaceApi struct {
}

func (r *NamespaceApi) GetNamespaceList(ctx *gin.Context) {

	todo := context.TODO()

	list, err := global.KubeConfigSet.CoreV1().Namespaces().List(todo, metav1.ListOptions{})
	if err != nil {
		response.FailWithMessage(ctx, err.Error())
		return
	}
	namespaceList := make([]namespace_res.Namespace, 0)

	for _, item := range list.Items {
		namespaceList = append(namespaceList, namespace_res.Namespace{
			Name:              item.Name,
			CreationTimestamp: item.CreationTimestamp.Unix(),
			Status:            string(item.Status.Phase),
		})
	}

	//fmt.Println(list.Items)

	response.SuccessWithDetailed(ctx, "获取成功", namespaceList)
}
