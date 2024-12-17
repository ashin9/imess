package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息通道
	MsgChan chan string
}

// 创建 Server 接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		MsgChan:   make(chan string),
	}
	return server
}

// 监听广播消息通道, 一旦有消息就广播给全部在线用户的消息通道
func (s *Server) ListenMessager() {
	for {
		msg := <-s.MsgChan
		// 将消息发送给全部用户
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.MsgChan <- msg
		}
		s.mapLock.Unlock()
	}
}

// 发送消息至广播消息通道
func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	s.MsgChan <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {
	// 根据连接新建用户
	user := NewUser(conn)

	// 用户上线, 添加到 OnlineMap 中
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	// 广播当前用户上线消息
	s.BroadCast(user, "已上线")
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

	// listen msgChan
	go s.ListenMessager()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listner accept err:", err)
			continue
		}
		// handler
		go s.Handler(conn)
	}
}
