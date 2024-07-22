package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	pb "t5/server/proto"
	"testing"
)

func TestClient(t *testing.T) {

	file, _ := credentials.NewClientTLSFromFile("D:\\Git\\go_note\\0project\\t5\\key\\test.pem", "*.tgo.com")

	//连接到 server端，此处禁用安全传输(没有加密和验证)
	//client, err := grpc.NewClient("127.0.0.1:9090",
	//	grpc.WithTransportCredentials(insecure.NewCredentials()))
	//
	client, err := grpc.NewClient("127.0.0.1:9090",
		grpc.WithTransportCredentials(file))

	if err != nil {
		log.Fatalf("未连接： %v", err)
	}
	defer client.Close()

	greeterClient := pb.NewGreeterClient(client)
	hello, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{Name: "01"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("返回值 %v", hello)

}
