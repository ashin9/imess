package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn, s)

	// 用户上线, 添加到 OnlineMap 中
	user.Online()

	// 监听用户活跃
	isAlive := make(chan bool)

	// 接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}

			// 提取用户发送的消息, 去除 '\n'
			msg := string(buf[:n-1])

			// 将得到的消息广播
			user.DoMessage(msg)

			// 更新用户活跃状态
			isAlive <- true
		}
	}()

	// 当前 handler 阻塞
	for {
		select {
		case <-isAlive:
			// 用户活跃, 重置定时器
			// 什么不做, 激活 select, 执行下面的定时器
		case <-time.After(time.Second * 60):
			// 超时, 强制关闭 User

			user.sendMsg("您已超时, 被踢下线")
			// 销毁用的资源
			close(user.MsgChan)
			// 关闭连接
			conn.Close()

			// 退出当前 handler
			return
		}
	}
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
