package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建 Server 接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (s *Server) handler(conn net.Conn) {
	// 连接当前业务
	fmt.Println("连接建立成功")
}

// 启动 Server 接口
func (s *Server) serve() {
	// listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen err:", err)
		return
	}
	// close listen
	defer listener.Close()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listner accept err:", err)
			continue
		}
		// handler
		go s.handler(conn)
	}
}
