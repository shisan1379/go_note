package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
	"user_growth/conf"
	"user_growth/dbhelper"
	"user_growth/pb"
	"user_growth/ugserver"
)

func initDb() {
	// default UTC time location
	time.Local = time.UTC
	// Load global config
	conf.LoadConfigs()
	// Initialize db
	dbhelper.InitDb()
}

func main() {
	// 初始化数据库实例
	initDb()

	listen, err := net.Listen("tcp", ":80")

	if err != nil {
		log.Fatal("监听失败", err)
	}
	s := grpc.NewServer()
	//注册服务
	pb.RegisterUserCoinServer(s, &ugserver.UgCoinServer{})   //用户积分
	pb.RegisterUserGradeServer(s, &ugserver.UgGradeServer{}) //用户等级

	reflection.Register(s)

	//启动服务
	log.Printf("服务启动中 %v\n", listen.Addr())
	err = s.Serve(listen)
	if err != nil {
		log.Fatal("服务启动失败：%v\n", err)
	}
}
