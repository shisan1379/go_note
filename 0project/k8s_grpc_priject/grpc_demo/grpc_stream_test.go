package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"testing"
)

func (c *RpcServer) ClientStream(stream service.Greeter_ClientStreamServer) error {
	count := 0
	for {
		//源源不断的去接收客户端发来的信息
		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println("服务端接收到的流", recv.Msg, count)
		count++
		if count > 10 {
			rsp := &service.HelloReply{Msg: "1"}
			err := stream.SendAndClose(rsp)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

// 服务端
func TestStreamServer(t *testing.T) {
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

// 客户端-使用 NewClient 方式创建 - 不验证
func TestStreamInsecureClient(t *testing.T) {
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
