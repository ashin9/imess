package main

import "net"

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	conn    net.Conn
}

// 创建用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		MsgChan: make(chan string),
		conn:    conn,
	}
	// 监听消息通道, 一有消息就发送给对应的客户端
	go user.ListenMessage()

	return user
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.MsgChan
		u.conn.Write([]byte(msg + "\n"))
	}
}
