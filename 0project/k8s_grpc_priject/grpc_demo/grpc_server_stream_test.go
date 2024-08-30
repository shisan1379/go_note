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

func (c *RpcServer) ServerStream(req *service.HelloRequest, stream service.Greeter_ServerStreamServer) error {
	for i := 0; i < 10; i++ {
		err := stream.Send(&service.HelloReply{
			Msg: strconv.Itoa(i),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
func (c *RpcServer) TwoStream(stream service.Greeter_TwoStreamServer) error {
	var i int
	for {
		i++
		recv, err := stream.Recv()
		if err != nil {
			return nil
		}
		fmt.Println("服务端收到客户端的消息", recv.Msg)
		time.Sleep(time.Second)
		rsp := &service.HelloReply{Msg: strconv.Itoa(i)}
		err = stream.Send(rsp)
		if err != nil {
			return nil
		}
	}
}

// 服务端
func TestStreamServerServer(t *testing.T) {
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
func TestStreamServerInsecureClient(t *testing.T) {
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

	stream, err := greeterClient.ServerStream(context.Background(), &service.HelloRequest{
		Msg: "Hello World",
	})
	if err != nil {
		log.Fatalf("获取服务端流失败 %v", err)
	}
	for i := 0; i < 10; i++ {

		recv, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端数据接收完成")
				err := stream.CloseSend()
				if err != nil {
					log.Fatal(err)
				}
				break
			}
			log.Fatal(err)
		}
		fmt.Println("服务端接收到的流", recv.Msg, i)

	}

	towSend(greeterClient)

}

func towSend(prodClient service.GreeterClient) {
	stream, err := prodClient.TwoStream(context.Background())

	var i int
	for {
		i++
		request := &service.HelloRequest{
			Msg: strconv.Itoa(i),
		}
		err = stream.Send(request)
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
		recv, err := stream.Recv()
		if err != nil {
			log.Fatal(err)
		}
		//websocket
		fmt.Println("客户端收到的流信息", recv.Msg)
	}
}
