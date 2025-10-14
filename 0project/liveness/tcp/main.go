package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// handleConnection 处理单个客户端的连接与消息交互
func handleConnection(conn net.Conn) {
	// 确保连接在函数结束后关闭
	defer func() {
		_ = conn.Close()
		fmt.Printf("客户端 [%s] 已断开连接\n", conn.RemoteAddr().String())
	}()

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("新客户端连接: [%s]\n", clientAddr)

	// 创建带缓冲的读取器，按行读取客户端消息
	reader := bufio.NewReader(conn)
	for {
		// 读取客户端发送的消息（以换行符 '\n' 为结束标志）
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取客户端 [%s] 消息失败: %v\n", clientAddr, err)
			return
		}

		// 去除消息两端的空白字符（包括换行符、空格等）
		cleanedMsg := strings.TrimSpace(msg)
		fmt.Printf("收到 [%s] 的消息: %q\n", clientAddr, cleanedMsg)

		// 根据消息内容生成响应
		var response string
		if cleanedMsg == "who" {
			response = "我\n" // 收到 "who" 时返回 "我"（加换行符方便客户端按行读取）
		} else {
			response = cleanedMsg + "\n" // 其他消息返回原值
		}

		// 发送响应给客户端
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("向客户端 [%s] 发送响应失败: %v\n", clientAddr, err)
			return
		}
	}
}

func main() {
	// 监听本地 8888 端口（TCP 协议）
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Printf("服务端启动失败（监听 8888 端口）: %v\n", err)
		return
	}
	// 确保监听器在程序结束后关闭
	defer func() {
		_ = listener.Close()
		fmt.Println("服务端已停止")
	}()

	fmt.Println("TCP 服务端已启动，监听端口 8888...")
	fmt.Println("等待客户端连接...")

	// 循环接受客户端连接（阻塞式）
	for {
		// 接受一个新的客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("接受客户端连接失败: %v\n", err)
			continue
		}

		// 启动协程处理该客户端的交互（支持并发）
		go handleConnection(conn)
	}
}
