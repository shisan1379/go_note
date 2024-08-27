package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"k8s_grpc_priject/grpc_demo/service"
)

func main() {
	hello := &service.HelloRequest{
		Msg:     "hello world",
		Address: []string{"1", "2", "3"},
	}
	// 序列化
	marshal, err := proto.Marshal(hello)
	if err != nil {
		panic(err)
	}
	//反序列化
	newHello := &service.HelloRequest{}
	err = proto.Unmarshal(marshal, newHello)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", newHello)
}
