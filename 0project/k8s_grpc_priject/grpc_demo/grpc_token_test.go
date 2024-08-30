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
)

// token 认证

func TestTokenServer(t *testing.T) {
	//拦截器
	var authInterceptor grpc.UnaryServerInterceptor
	authInterceptor = func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		//拦截普通方法请求，验证 Token
		err = Auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}
	//启动服务时注册拦截器
	server := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))
	service.RegisterGreeterServer(server, &RpcServer{})

	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}
	err = server.Serve(listener)
	if err != nil {
		log.Fatal("服务、启动失败", err)
	}
	fmt.Println("启动成功")
}
func Auth(ctx context.Context) error {
	//取出值
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	var user string
	var password string

	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if user != "admin" || password != "admin" {
		return status.Errorf(codes.Unauthenticated, "token不合法")
	}
	return nil
}

// 客户端 实现 验证器 - PerRPCCredentials
type Authentication struct {
	User     string
	Password string
}

// 实现获取数据方法
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (
	map[string]string, error,
) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

// 是否需要 tls 的安全性
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func TestTokenClient(t *testing.T) {
	user := &Authentication{
		User:     "admin",
		Password: "admin",
	}

	//设置出站凭据
	//grpc.WithPerRPCCredentials(user)

	client, err := grpc.NewClient(":8002", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(user))

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
