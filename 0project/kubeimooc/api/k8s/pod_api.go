package k8s

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubeimooc/global"
	"net/http"
)

type PodApi struct {
}

func (r *PodApi) GetPodList(c *gin.Context) {

	todo := context.TODO()
	// 这里为空就是查询namespace
	list, err := global.KubeConfigSet.CoreV1().Pods("").List(todo, metav1.ListOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, item := range list.Items {
		fmt.Println(item.Namespace, item.Name)
	}

	c.JSON(http.StatusOK, gin.H{})
}
