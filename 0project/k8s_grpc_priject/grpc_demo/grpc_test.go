package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/anypb"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

type RpcServer struct {
	service.UnimplementedGreeterServer
}

func (c *RpcServer) SayHello(ctx context.Context, req *service.HelloRequest) (*service.HelloReply, error) {
	log.Printf("请求值： %v", req)
	a, _ := anypb.New(&service.DataMsg{Data: "data is me"})
	return &service.HelloReply{
		Msg:  "hello",
		Data: a,
	}, nil
}

// 服务端
func TestServer(t *testing.T) {
	listen, _ := net.Listen("tcp", ":9090")
	rpcServer := grpc.NewServer()
	service.RegisterGreeterServer(rpcServer, &RpcServer{})

	err := rpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		panic(err)
	} else {
		fmt.Println("RpcServer started")
	}

}

// 客户端-使用 Dial 方式创建 - 不验证
func TestInsecureClientDial(t *testing.T) {

	client, err := grpc.Dial("127.0.0.1:9090",
		grpc.WithInsecure(),           // 这个选项告诉 gRPC 客户端忽略 TLS 证书验证
		grpc.WithBlock(),              // 这个选项会让 Dial 阻塞，直到连接建立或发生错误
		grpc.WithTimeout(time.Second)) // 设置连接超时时间为 1 秒
	if err != nil {
		log.Fatalf("未连接： %v", err)
	}
	defer client.Close()

	greeterClient := service.NewGreeterClient(client)
	hello, err := greeterClient.SayHello(context.Background(), &service.HelloRequest{Msg: "01"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("返回值 %v", hello)
}

// 客户端-使用 NewClient 方式创建 - 不验证
func TestInsecureClient(t *testing.T) {
	// 创建一个不安全的客户端凭据，这通常用于测试环境，不建议在生产环境中使用
	cred := insecure.NewCredentials()
	// 使用上述凭据配置gRPC的传输凭据
	transportCredentials := grpc.WithTransportCredentials(cred)

	//尝试使用配置的凭据连接到gRPC服务器
	client, err := grpc.NewClient("127.0.0.1:9090", transportCredentials)

	if err != nil {
		log.Fatalf("未连接： %v", err) //如果错误则退出
	}
	defer client.Close() //最后关闭 连接

	//创建一个Greeter服务的客户端
	greeterClient := service.NewGreeterClient(client)
	background := context.Background()
	// 调用Greeter服务的SayHello方法，发送请求并等待响应
	hello, err := greeterClient.SayHello(background,
		&service.HelloRequest{
			Msg:  "01",
			User: &service.User{Name: "123"},
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	var dd service.DataMsg
	err = hello.Data.UnmarshalTo(&dd)

	fmt.Printf("返回值 %s , %s \n", hello.Msg, dd.Data)
}

type userClient struct {
	clients []RpcServer
	index   int64
}

// 简单池化
func (r *userClient) Get() RpcServer {
	//1. index + 1，相当于每轮询使用 池中的资源
	index := atomic.AddInt64(&r.index, 1)
	i := int(index) % len(r.clients)
	return r.clients[i]
}

func NewClientPool(size int) {
	var cs []RpcServer
	for i := 0; i < size; i++ {
		server := RpcServer{}
		cs = append(cs, server)
	}
}

func (r *userClient) Release() {

}
