package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"user_growth/pb"
	"user_growth/ugserver"
)

func main() {
	listen, err := net.Listen("tcp", ":80")

	if err != nil {
		log.Fatal("监听失败", err)
	}
	s := grpc.NewServer()
	//注册服务
	pb.RegisterUserCoinServer(s, &ugserver.UgCoinServer{})    //用户积分
	pb.RegisterUserGradeServer(s, &ugserver.UgGrowthServer{}) //用户等级

	//启动服务
	log.Printf("服务启动中 %v\n", listen.Addr())
	err = s.Serve(listen)
	if err != nil {
		log.Fatal("服务启动失败：%v\n", err)
	}

}
