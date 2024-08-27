package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"testing"
	"time"
)

type server struct {
	service.UnimplementedGreeterServer
}

func (c *server) SayHello(ctx context.Context, req *service.HelloRequest) (*service.HelloReply, error) {
	fmt.Printf("请求值 %v", *req)
	return &service.HelloReply{
		Msg: "hello",
	}, nil
}

// 服务端
func TestServer(t *testing.T) {
	listen, _ := net.Listen("tcp", ":9090")
	rpcServer := grpc.NewServer()
	service.RegisterGreeterServer(rpcServer, &server{})

	err := rpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		panic(err)
	} else {
		fmt.Println("server started")
	}

}

// 客户端
func TestClient(t *testing.T) {
	//client, err := grpc.NewClient("127.0.0.1:9090")
	client, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Second))
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
