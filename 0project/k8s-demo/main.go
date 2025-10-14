package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	router := gin.New()
	router.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, getClusterIp()+" "+getDNS()+" : hello")
	})

	router.GET("/healthz", func(ctx *gin.Context) {
		fmt.Println("healthz")
		ctx.String(http.StatusOK, "ok")
	})

	router.GET("/notify", func(ctx *gin.Context) {
		fmt.Println("notify")
		ctx.String(http.StatusOK, "ok")
	})

	// 为http/2配置参数
	h2Handler := h2c.NewHandler(router, &http2.Server{})
	// 配置http服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2Handler,
	}
	// 启动http服务
	server.ListenAndServe()
}
func getClusterIp() string {
	// 从环境变量中获取 ClusterIP 和端口
	serviceHost := os.Getenv("MY_SERVICE_SERVICE_HOST")
	servicePort := os.Getenv("MY_SERVICE_SERVICE_PORT")

	if serviceHost == "" || servicePort == "" {
		fmt.Println("无法获取 Service 的 ClusterIP 或端口")
		return "0---"
	}

	fmt.Printf("Service 的 ClusterIP: %s\n", serviceHost)
	fmt.Printf("Service 的端口: %s\n", servicePort)
	return serviceHost
}
func getClusterIp2() string {
	// 加载 kubeconfig 文件
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, "..kube", "config")
	} else {
		log.Fatal("无法找到 kubeconfig 文件")
	}

	// 创建配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("无法创建配置: %v", err)
	}

	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("无法创建客户端: %v", err)
	}

	// 获取 Service 信息
	serviceName := "my-service"
	namespace := "default"
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("无法获取 Service: %v", err)
	}

	// 输出 ClusterIP
	fmt.Printf("Service %s 的 ClusterIP: %s\n", serviceName, service.Spec.ClusterIP)

	return service.Spec.ClusterIP
}
func getDNS() string {
	serviceName := "myapp.default.svc.cluster.local"
	addrs, err := net.LookupHost(serviceName)
	if err != nil {
		log.Fatalf("Failed to resolve %s: %v", serviceName, err)
	}
	clusterIP := addrs[0]
	return clusterIP
}
