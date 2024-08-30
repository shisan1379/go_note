// 客户端流 demo

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
	"strconv"
	"testing"
	"time"
)

// 客户端流处理服务实现
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

	stream, err := greeterClient.ClientStream(context.Background())
	if err != nil {
		log.Fatal("获取流出错", err)
	}
	// 创建一个基于 struct 的 channel，容量为1
	rsp := make(chan struct{}, 1)
	// 10 发送消息
	go send10TimesRequest(stream, rsp)
	// 当发送10次后，等待返回值
	select {
	case <-rsp:
		//关闭并接收返回值
		recv, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatal(err)
		}
		stock := recv.Msg
		fmt.Println("客户端收到响应：", stock)
	}
}
func send10TimesRequest(stream service.Greeter_ClientStreamClient, rsp chan struct{}) {
	count := 0
	var i int
	for {
		i++
		request := &service.HelloRequest{
			Msg: strconv.Itoa(i),
		}
		// 基于流发送消息
		err := stream.SendMsg(request)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		count++
		if count > 10 {
			rsp <- struct{}{}
			break
		}
	}
}
