package main

import (
	"net"
	"strings"
)

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
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 修改用户名, 消息格式 "rename|张三"
		newName := strings.Split(msg, "|")[1]

		// 判断 name 是否已经占用
		_, ok := u.Server.OnlineMap[newName]
		if ok {
			u.sendMsg("当前用户名已经被占用\n")
		} else {
			u.Server.mapLock.Lock()
			delete(u.Server.OnlineMap, u.Name)
			u.Server.OnlineMap[newName] = u
			u.Server.mapLock.Unlock()

			u.Name = newName
			u.sendMsg("您已经更新用户名:" + u.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		// 获取对方用户名
		toName := strings.Split(msg, "|")[1]
		if toName == "" {
			u.sendMsg("消息格式不正确, 请使用 \"to|张三|你好\" 格式. \n")
			return
		}
		// 根据用户名查找对方 User 对象
		toUser, ok := u.Server.OnlineMap[toName]
		if !ok {
			u.sendMsg("您输入的用户名不存在. \n")
			return
		}
		// 获取发送消息
		toMsg := strings.Split(msg, "|")[2]
		if toMsg == "" {
			u.sendMsg("无消息内容, 请重新发送\n")
			return
		}
		// 发送消息
		toUser.sendMsg(u.Name + "对您说: " + toMsg + "\n")

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
