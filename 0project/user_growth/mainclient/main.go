package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
	"user_growth/pb"
)

func main() {
	//连接到服务端
	add := flag.String("addr", "127.0.0.1:8080", "server listen address")

	// 创建一个不安全的客户端凭据，这通常用于测试环境，不建议在生产环境中使用
	cred := insecure.NewCredentials()
	// 使用上述凭据配置gRPC的传输凭据
	transportCredentials := grpc.WithTransportCredentials(cred)

	//尝试使用配置的凭据连接到gRPC服务器
	client, err := grpc.NewClient(*add, transportCredentials)
	if err != nil {
		log.Fatal("客户端连接服务器失败，%v", err)
	}
	defer client.Close()

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*1)
	defer cancelFunc()

	//新建客户端对象
	coinClient := pb.NewUserCoinClient(client)
	gradeClient := pb.NewUserGradeClient(client)

	//测试1
	r1, err := coinClient.ListTask(ctx, &pb.ListTaskRequest{})
	if err != nil {
		log.Printf("coin ListTask err:%v\n", err)
	} else {
		log.Printf("list task result:%v\n", r1)
	}

	r2, err := gradeClient.ListGrades(ctx, &pb.ListGradesRequest{})
	if err != nil {
		log.Printf("coin ListTask err:%v\n", err)
	} else {
		log.Printf("list task result:%v\n", r2)
	}
}
