package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	pb "t5/server/proto"
	"testing"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (c *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{
		Message: "123",
	}, nil
}
func TestServer(t *testing.T) {

	file, _ := credentials.NewServerTLSFromFile("D:\\Git\\go_note\\0project\\t5\\key\\test.pem",
		"D:\\Git\\go_note\\0project\\t5\\key\\test.key")
	//开启端口
	listen, _ := net.Listen("tcp", ":9090")
	//创建GRPC服务
	//grpcServer := grpc.NewServer()
	grpcServer := grpc.NewServer(grpc.Creds(file))
	//在grpc中注册服务
	pb.RegisterGreeterServer(grpcServer, &server{})

	//启动服务
	err := grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("启动服务失败")
	}
}
