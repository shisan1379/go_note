package global

import (
	"k8s.io/client-go/kubernetes"
	"kubeimooc/config"
)

var (
	//配置文件
	CONF          config.Server
	KubeConfigSet *kubernetes.Clientset
)
