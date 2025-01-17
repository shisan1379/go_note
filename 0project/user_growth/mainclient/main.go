package main

import (
	"context"
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"sync"
	"time"
	"user_growth/pb"
)

var connPool = sync.Pool{
	New: func() any {
		// 连接到服务
		addr := flag.String("addr", "127.0.0.1:80", "server listen address")

		opts := []grpc.DialOption{
			// 创建一个不安全的客户端凭据，这通常用于测试环境，不建议在生产环境中使用
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			//gRPC客户端的读写缓冲区设置
			grpc.WithWriteBufferSize(1024 * 1024 * 1), // 默认32KB
			grpc.WithReadBufferSize(1024 * 1024 * 1),  // 默认32KB,
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Minute, // 指定客户端多久发送一个保活探测消息给服务端
				Timeout:             10 * time.Second, // 指定等待服务端响应保活探测消息的最长时间
				PermitWithoutStream: false,            // 当没有活动的RPC流时，是否允许发送保活探测消息
			}),
		}

		//尝试使用配置的凭据连接到gRPC服务器
		client, err := grpc.NewClient(*addr, opts...)

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return client
	},
}

func GetConn() *grpc.ClientConn {
	return connPool.Get().(*grpc.ClientConn)
}

func CloseConn(conn *grpc.ClientConn) {
	connPool.Put(conn)
}

func main() {
	conn := GetConn()
	if conn != nil {
		defer CloseConn(conn)
	} else {
		log.Fatalf("connection nil")
	}
	// 请求服务
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 新建客户端
	cCoin := pb.NewUserCoinClient(conn)
	cGrade := pb.NewUserGradeClient(conn)

	// 测试1：UserCoinServer.ListTasks
	r1, err1 := cCoin.ListTask(ctx, &pb.ListTaskRequest{})
	if err1 != nil {
		log.Printf("cCoin.ListTasks error=%v\n", err1)
	} else {
		log.Printf("cCoin.ListTasks: %+v\n", r1.GetDatalist())
	}
	// 测试2：UserGradeServer.ListGrades
	r2, err2 := cGrade.ListGrades(ctx, &pb.ListGradesRequest{})
	if err2 != nil {
		log.Printf("cGrade.ListGrades error=%v\n", err2)
	} else {
		log.Printf("cGrade.ListGrades: %+v\n", r2.GetDatalist())
	}
	// 测试3：修改积分
	r3, err3 := cCoin.UserCoinChange(ctx, &pb.UserCoinChangeRequest{
		Uid:  0,
		Task: "abc",
		Coin: 0,
	})
	if err3 != nil {
		log.Printf("cCoin.UserCoinChange error=%v\n", err3)
	} else {
		log.Printf("cCoin.UserCoinChange: %+v\n", r3.GetUser())
	}
}
