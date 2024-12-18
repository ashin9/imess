package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	Mode       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		Mode:       1,
	}

	// 连接 Server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net dial err:", err)
		return nil
	}

	client.Conn = conn

	// 返回对象
	return client
}

func (client *Client) menu() bool {
	var mode int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&mode)

	if mode >= 0 && mode <= 3 {
		client.Mode = mode
		return true
	} else {
		fmt.Println(">>>>请输入合法数字<<<<")
		return false
	}
}

func (client *Client) revResponse() {
	io.Copy(os.Stdout, client.Conn)

	// 等价于
	// for {
	// 	buf := make([]byte, 4096)
	// 	client.Conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func (client *Client) PublicChat() {
	var puclicMsg string

	fmt.Println(">>>> 请输入聊天内容, exit退出.")
	fmt.Scanln(&puclicMsg)

	for puclicMsg != "exit" {
		if len(puclicMsg) != 0 {
			_, err := client.Conn.Write([]byte(puclicMsg + "\n"))
			if err != nil {
				fmt.Println("conn write err:", err)
				return
			}
		}

		puclicMsg = ""
		fmt.Println(">>>> 请输入聊天内容, exit退出.")
		fmt.Scanln(&puclicMsg)
	}

}

func (client *Client) SelectUsers() {
	sendMsg := "who"
	_, err := client.Conn.Write([]byte(sendMsg + "\n"))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var privateName string

	client.SelectUsers()
	fmt.Println(">>>> 请输入要私聊的用户(exit退出):")
	fmt.Scanln(&privateName)

	for privateName != "exit" {

		var privateMsg string
		fmt.Println(">>>> 请输入要私聊的消息(exit退出):")
		fmt.Scanln(&privateMsg)

		for privateMsg != "exit" {
			if len(privateMsg) != 0 {
				sendMsg := "to|" + privateName + "|" + privateMsg + "\n"
				_, err := client.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err:", err)
					break
				}
			}

			privateMsg = ""
			fmt.Println(">>>> 请输入要私聊的消息(exit退出):")
			fmt.Scanln(&privateMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>> 请输入要私聊的用户(exit退出):")
		fmt.Scanln(&privateName)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>>请输入更改后的用户名:")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}

func (client *Client) Run() {
	for client.Mode != 0 {
		for !client.menu() {
		}

		// 根据不同业务选择不同模式
		switch client.Mode {
		case 1:
			// fmt.Println("公聊模式...")
			client.PublicChat()
			break
		case 2:
			// fmt.Println("私聊模式...")
			client.PrivateChat()
			break
		case 3:
			// fmt.Println("更新用户名...")
			client.UpdateName()
			break
		}

	}
}
