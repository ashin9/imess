package main

import (
	"flag"
	"fmt"
)

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器 IP (默认为 127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器 Port (默认为 8888)")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>> 连接服务器失败...")
		return
	}
	fmt.Println(">>>> 连接服务器成功...")

	// 启动 goroutine 监听 server 端响应的消息
	go client.revResponse()
	// 启动客户端业务, 发送消息
	client.Run()
}
