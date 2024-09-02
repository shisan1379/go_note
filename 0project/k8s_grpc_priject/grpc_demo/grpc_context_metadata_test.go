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
	// UnaryServerInterceptor 是一元 RPC 调用的拦截器

	interceptor := grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

			if info.FullMethod == service.Greeter_SayHello_FullMethodName {
				fmt.Println("当前请求方法：", service.Greeter_SayHello_FullMethodName)
				//登录
				err = Auth(ctx)
				if err != nil {
					return
				}
			}

			//将  信息传递给 context ， 便于后续拦截器使用
			ctx = context.WithValue(ctx, "login", true)

			return handler(ctx, req)
		},
		//
		func(ctx context.Context, //上下文
			req any, //具体的请求
			info *grpc.UnaryServerInfo, // grpc server
			handler grpc.UnaryHandler, //真正实现服务的方法
		) (resp any, err error) {
			//记录时间
			now := time.Now()
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				//未认证
				return ctx, status.Errorf(codes.Unauthenticated, "未提供认证数据")
			}

			value := ctx.Value("login")
			fmt.Println("login", value)

			var r any
			tokens := md.Get("token")
			keys := md.Get("key")
			if len(keys) > 0 {
				r = keys[0]
				fmt.Println("key", r)
			} else {
				fmt.Println("未传入 keys")
			}

			if len(tokens) > 0 {
				tk := tokens[0]
				if tk == "123" {
					// 继续处理请求
					newMd := metadata.Pairs("token", "server token  is 123")
					grpc.SendHeader(ctx, newMd) // 确保在发送任何响应之前调用

					//md := metadata.Pairs("token", "value1", "key2", "value2")
					//outCtx := metadata.NewOutgoingContext(ctx, md)

					r, err = handler(ctx, req)
				} else {
					return ctx, status.Errorf(codes.Unauthenticated, "未提供认证数据1")
				}
			} else {

				return ctx, status.Errorf(codes.Unauthenticated, "未提供认证数据2")
			}

			//打印用时
			fmt.Printf("用时 %v\n", time.Since(now))

			return r, err
		},
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

	user := &Authentication{
		User:     "admin",
		Password: "admin",
	}

	// 创建一个不安全的客户端凭据，这通常用于测试环境，不建议在生产环境中使用
	cred := insecure.NewCredentials()
	// 使用上述凭据配置gRPC的传输凭据
	transportCredentials := grpc.WithTransportCredentials(cred)

	// 使用客户端拦截器进行客户端统一的metadata处理
	// 拦截器按照注册顺序执行
	interceptor := grpc.WithChainUnaryInterceptor(func(ctx context.Context,
		method string, req, reply any, c *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.Pairs("token", "123")

		newCtx := metadata.NewOutgoingContext(ctx, md)

		newCtx = metadata.AppendToOutgoingContext(newCtx, "key", "k1")

		// 使用新的context调用原始的RPC方法, 并获取服务端设置的 header
		var header, trailer metadata.MD
		err := invoker(newCtx, method, req, reply, c, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			// 处理错误（如果需要）
			return err
		}

		fmt.Println(method+":server ts", header.Get("token"))

		//trailer
		// 如果需要，可以在这里进一步处理reply或进行日志记录
		return nil
	},
	)

	//尝试使用配置的凭据连接到gRPC服务器
	client, err := grpc.NewClient("127.0.0.1:9090", transportCredentials, interceptor, grpc.WithPerRPCCredentials(user))

	if err != nil {
		log.Fatalf("未连接： %v", err) //如果错误则退出
	}
	defer client.Close() //最后关闭 连接

	//创建一个Greeter服务的客户端
	greeterClient := service.NewGreeterClient(client)
	//
	//// 调用Greeter服务的SayHello方法，发送请求并等待响应
	//hello, err := greeterClient.SayHello(context.Background(),
	//	&service.HelloRequest{
	//		Msg:  "01",
	//		User: &service.User{Name: "123"},
	//	},
	//)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//var dd service.DataMsg
	//err = hello.Data.UnmarshalTo(&dd)
	//
	//fmt.Printf("返回值 %s , %s \n", hello.Msg, dd.Data)

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
