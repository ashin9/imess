package main

import "net"

type User struct {
	Name    string
	Addr    string
	MsgChan chan string
	Conn    net.Conn

	Server *Server
}

// 创建用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:    userAddr,
		Addr:    userAddr,
		MsgChan: make(chan string),
		Conn:    conn,
		Server:  server,
	}
	// 监听消息通道, 一有消息就发送给对应的客户端
	go user.ListenMessage()

	return user
}

// 用户上线业务
func (u *User) Online() {
	// 用户上线, 添加到 OnlineMap 中
	u.Server.mapLock.Lock()
	u.Server.OnlineMap[u.Name] = u
	u.Server.mapLock.Unlock()

	// 广播当前用户上线消息
	u.Server.BroadCast(u, "已上线")
}

// 用户下线业务
func (u *User) Offline() {
	// 用户下线, 删除 OnlineMap 中
	u.Server.mapLock.Lock()
	delete(u.Server.OnlineMap, u.Name)
	u.Server.mapLock.Unlock()

	// 广播当前用户下线消息
	u.Server.BroadCast(u, "已下线")
}

// 给当前用户客户端发送消息
func (u *User) sendMsg(msg string) {
	u.Conn.Write([]byte(msg))
}

// 用户处理消息
func (u *User) DoMessage(msg string) {
	// 查询在线用户
	if msg == "who" {
		u.Server.mapLock.Lock()
		for _, user := range u.Server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			u.sendMsg(onlineMsg)
		}
		u.Server.mapLock.Unlock()
	} else {
		u.Server.BroadCast(u, msg)
	}
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.MsgChan
		u.Conn.Write([]byte(msg + "\n"))
	}
}
