package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"k8s_grpc_priject/grpc_demo/service"
	"log"
	"net"
	"testing"
)

//双向 TLS

func TestTwoTlsServer(t *testing.T) {
	//添加证书
	//file, err2 := credentials.NewServerTLSFromFile("keys/mszlu.pem", "keys/mszlu.key")
	//if err2 != nil {
	//	log.Fatal("证书生成错误",err2)
	//}
	// 证书认证-双向认证
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, err := tls.LoadX509KeyPair("./keys/server.pem", "./keys/server.key")
	if err != nil {
		log.Fatal("证书读取错误", err)
	}
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("./keys/ca.crt")
	if err != nil {
		log.Fatal("ca证书读取错误", err)
	}
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ClientAuth: tls.RequireAndVerifyClientCert,
		// 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
		ClientCAs: certPool,
	})

	rpcServer := grpc.NewServer(grpc.Creds(creds))

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
func TestTwoTlsClient(t *testing.T) {
	//file, err2 := credentials.NewClientTLSFromFile("client/keys/mszlu.pem", "*.mszlu.com")
	//if err2 != nil {
	//	log.Fatal("证书错误",err2)
	//}
	// 证书认证-双向认证
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, _ := tls.LoadX509KeyPair("./keys/client.pem", "client/keys/client.key")
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("./keys/ca.crt")
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)
	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		// 要求必须校验客户端的证书。可以根据实际情况选用以下参数
		ServerName: "*.mszlu.com",
		RootCAs:    certPool,
	})

	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(creds))

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
