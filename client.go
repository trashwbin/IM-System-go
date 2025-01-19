package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

// NewClient 创建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建一个客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}

	client.conn = conn

	// 返回一个客户端
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败...")
		return
	}

	fmt.Println(">>>>> 连接服务器成功...")

	//select {}
}
