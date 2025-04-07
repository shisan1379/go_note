package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeimooc/global"
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
	fmt.Println(list.Items)

	response.Success(ctx)
}
