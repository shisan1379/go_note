package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"os"
	"testing"
)

//双向 TLS

func TestTwo2TlsServer(t *testing.T) {
	// 加载CA证书
	caCert, err := os.ReadFile("./keys2/ca.pem")
	if err != nil {
		log.Fatalf("ca证书加载失败: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 加载服务器证书和密钥
	serverCert, err := tls.LoadX509KeyPair("./keys2/server.pem", "./keys2/server.key")
	if err != nil {
		log.Fatalf("服务端证书加载失败: %v", err)
	}

	// 配置TLS
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	})

	rpcServer := grpc.NewServer(grpc.Creds(creds))
	service.RegisterGreeterServer(rpcServer, &RpcServer{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("启动监听出错", err)
	}
	err = rpcServer.Serve(listener)
	if err != nil {
		log.Fatal("启动服务出错", err)
	}
	fmt.Println("启动grpc服务端成功")
}
func TestTwo2TlsClient(t *testing.T) {

	// 加载CA证书
	caCert, err := os.ReadFile("./keys2/ca.pem")
	if err != nil {
		log.Fatalf("ca证书加载失败: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 加载客户端证书和密钥
	clientCert, err := tls.LoadX509KeyPair("./keys2/client.pem", "./keys2/client.key")
	if err != nil {
		log.Fatalf("客户端证书加载失败: %v", err)
	}

	// 配置TLS
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{clientCert},
		ServerName:   "*.mszlu.com",
		RootCAs:      caCertPool,
	})

	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
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
