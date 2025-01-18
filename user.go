package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

// 创建一个用户的API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前user channel消息的goroutine
	go user.ListenMessage()

	return user
}

// 监听当前User channel的方法，一旦有消息，就直接发送给对端客户端
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}

// 用户上线业务
func (user *User) Online() {
	// 用户上线，将用户加入到onlineMap中
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	// 广播当前用户上线消息
	//server.BroadCast(user, "已上线")
	user.server.BroadCast(user, "has online")
}

// 用户下线业务
func (user *User) Offline() {
	// 用户下线，将用户从onlineMap中删除
	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.mapLock.Unlock()

	// 广播当前用户下线消息
	user.server.BroadCast(user, "has offline")
}

// 用户处理消息的业务
func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// 查询当前在线用户
		user.server.mapLock.Lock()
		for _, onlineUser := range user.server.OnlineMap {
			onlineMsg := "[" + onlineUser.Addr + "]" + onlineUser.Name + ": online...\n"
			user.SendMessage(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		// 消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		// 判断name是否存在
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMessage("当前用户名已经被使用\n")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.SendMessage("您已经更新用户名：" + user.Name + "\n")
		}

	} else {
		// 通过广播消息
		user.server.BroadCast(user, msg)
	}
}

// 给当前用户的对应客户端发送消息
func (user *User) SendMessage(msg string) {
	user.conn.Write([]byte(msg))
}
