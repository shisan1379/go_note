package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"testing"
)

//单向 TLS

// 服务端
func TestTlsServer(t *testing.T) {
	//添加证书
	file, err2 := credentials.NewServerTLSFromFile("./keys/server.pem", "./keys/server.key")
	if err2 != nil {
		log.Fatal("证书生成错误", err2)
	}
	rpcServer := grpc.NewServer(grpc.Creds(file))

	service.RegisterGreeterServer(rpcServer, &server{})

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatal("启动监听出错", err)
	}
	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("启动服务出错", err)
	}
	fmt.Println("启动grpc服务端成功")
}

// 客户端
func TestClientDial(t *testing.T) {
	file, err2 := credentials.NewClientTLSFromFile("./keys/server.pem", "*.mszlu.com")
	if err2 != nil {
		log.Fatal("证书错误", err2)
	}
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(file))

	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	defer conn.Close()

	prodClient := service.NewGreeterClient(conn)

	request := &service.HelloRequest{
		Msg: "123",
	}
	stockResponse, err := prodClient.SayHello(context.Background(), request)
	if err != nil {
		log.Fatal("查询出错", err)
	}
	fmt.Println("查询成功", stockResponse)
}

func TestClient(t *testing.T) {
	file, err2 := credentials.NewClientTLSFromFile("./keys/server.pem", "*.mszlu.com")
	if err2 != nil {
		log.Fatal("证书错误", err2)
	}
	conn, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(file))

	if err != nil {
		log.Fatal("服务端出错，连接不上", err)
	}
	defer conn.Close()

	prodClient := service.NewGreeterClient(conn)

	request := &service.HelloRequest{
		Msg: "123",
	}
	stockResponse, err := prodClient.SayHello(context.Background(), request)
	if err != nil {
		log.Fatal("查询出错", err)
	}
	fmt.Println("查询成功", stockResponse)
}
