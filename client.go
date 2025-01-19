package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int // 当前客户端模式
}

// NewClient 创建一个客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建一个客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
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

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口（默认是8888）")
}

func (client *Client) menu() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("1. 公聊模式\n2. 私聊模式\n3. 更新用户名\n0. 退出")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取输入时出错:", err)
		return false
	}

	input = strings.TrimSpace(input) // 去掉换行符和空白字符
	var flag int
	_, err = fmt.Sscanf(input, "%d", &flag) // 解析输入的数字
	if err != nil {
		fmt.Println(">>>>> 请输入合法范围内的数字...")
		return false
	}

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>> 请输入合法范围内的数字...")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
		case 2:
			// 私聊模式
			client.PrivateChat()
		case 3:
			// 更新用户名
			client.UpdateName()
		}
	}
}

// 查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}
func (client *Client) PrivateChat() {
	reader := bufio.NewReader(os.Stdin)

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")

	remoteName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取输入时出错:", err)
		return
	}
	remoteName = strings.TrimSpace(remoteName) // 去掉换行符和空白字符

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容, exit退出:")
		chatMsg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入时出错:", err)
			return
		}
		chatMsg = strings.TrimSpace(chatMsg) // 去掉换行符和空白字符

		for chatMsg != "exit" {
			// 消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			fmt.Println(">>>>请输入消息内容, exit退出:")
			chatMsg, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("读取输入时出错:", err)
				return
			}
			chatMsg = strings.TrimSpace(chatMsg) // 去掉换行符和空白字符
		}

		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
		remoteName, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入时出错:", err)
			return
		}
		remoteName = strings.TrimSpace(remoteName) // 去掉换行符和空白字符
	}
}

func (client *Client) PublicChat() {
	// 提示用户输入消息
	fmt.Println(">>>>> 请输入聊天内容，exit退出")

	reader := bufio.NewReader(os.Stdin)
	for {
		chatMsg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入时出错:", err)
			return
		}

		// 去掉所有空白字符
		chatMsg = strings.TrimSpace(chatMsg)

		if chatMsg == "exit" {
			break
		}

		// 发送给服务器
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}

		fmt.Println(">>>>> 请输入聊天内容，exit退出")
	}
}

func (client *Client) UpdateName() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(">>>>> 请输入用户名：")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取输入时出错:", err)
		return false
	}

	client.Name = strings.TrimSpace(input) // 去掉换行符和空白字符
	msg := "rename|" + client.Name + "\n"
	_, err = client.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

// DealResponse 处理server的响应
func (client *Client) DealResponse() {
	// 一旦client.conn有数据，就会拷贝到os.Stdout，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>> 连接服务器失败...")
		return
	}

	fmt.Println(">>>>> 连接服务器成功...")
	// 单独开启一个goroutine处理server的响应
	go client.DealResponse()

	client.Run()
}

/*
	func (client *Client) PublicChat() {
		//提示用户输入消息
		var chatMsg string

		fmt.Println(">>>>请输入聊天内容，exit退出.")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//发给服务器

			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入聊天内容，exit退出.")
			fmt.Scanln(&chatMsg)
		}

}
*/
/*
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容, exit退出:")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			//消息不为空则发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入消息内容, exit退出:")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象[用户名], exit退出:")
		fmt.Scanln(&remoteName)
	}
}
*/
