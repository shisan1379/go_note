// grpc 中 context 和 metadata 使用实例

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"testing"
	"time"
)

func TestContextMetadataServer(t *testing.T) {
	//单个拦截器
	//var authInterceptor grpc.UnaryServerInterceptor
	//grpc.UnaryInterceptor(authInterceptor)

	//链式拦截器
	interceptor := grpc.ChainUnaryInterceptor(
		func(ctx context.Context,       //上下文
			req any,                    //具体的请求
			info *grpc.UnaryServerInfo, // grpc server
			handler grpc.UnaryHandler,  //真正实现服务的方法
		) (resp any, err error) {
			//记录时间
			now := time.Now()
			log.Printf("1")
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				//未认证
				return ctx, status.Errorf(codes.Unauthenticated, "未提供认证数据")
			}

			var r any
			tokens := md.Get("token")
			if len(tokens) > 0 {
				tk := tokens[0]
				if tk != "123" {
					// 继续处理请求
					r, err = handler(ctx, req)
				} else {
					return ctx, status.Errorf(codes.Unauthenticated, "未提供认证数据")
				}
			} else {

				return ctx, status.Errorf(codes.Unauthenticated, "未提供认证数据")
			}

			//打印用时
			fmt.Printf("用时 %v\n", time.Since(now))

			return r, err
		},
		//func(ctx context.Context, //上下文
		//	req any, //具体的请求
		//	info *grpc.UnaryServerInfo, // grpc server
		//	handler grpc.UnaryHandler, //真正实现服务的方法
		//) (resp any, err error) {
		//	//记录时间
		//	log.Printf("2")
		//	// 继续处理请求
		//	a, err := handler(ctx, req)
		//
		//	return a, err
		//},
	)

	//启动服务时注册拦截器
	server := grpc.NewServer(interceptor)
	service.RegisterGreeterServer(server, &RpcServer{})

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}
	err = server.Serve(listener)
	if err != nil {
		log.Fatal("服务、启动失败", err)
	}
	fmt.Println("启动成功")
}
func TestContextMetadataInsecureClient(t *testing.T) {
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
	//md := metadata.MD{
	//	"token": []string{"1234"},
	//}
	md := metadata.Pairs("token", "1234")
	background = metadata.NewOutgoingContext(background, md)

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
