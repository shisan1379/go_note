package main

import (
	"fmt"
	"os"
	"runtime/trace"
)

func main() {

	//创建trace文件
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	//启动trace goroutine
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	//for i := 0; i < 10; i++ {
	//	go func() {
	//		fmt.Println(strconv.Itoa(i) + "goroutine")
	//	}()
	//}

	//main
	fmt.Println("Hello World")
}
