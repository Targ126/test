package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//用户表
	UserMap map[string]*User

	// map 锁
	mapLock sync.RWMutex

	// 管道
	Message chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:      ip,
		Port:    port,
		UserMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}

// Start server启动
func (server *Server) Start() {
	//监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	//关闭
	defer listener.Close()

	//启动广播的监听
	go server.ListenMessage()
	// 遍历轮询 检查是否连接成功
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		server.Handler(conn)
	}
}

// ListenMessage 监听Message中消息，并发送给所有在线用户
func (server *Server) ListenMessage() {
	for {
		msg := <-server.Message

		server.mapLock.Lock()
		for _, user := range server.UserMap {
			user.C <- msg
		}
		server.mapLock.Unlock()
	}
}

// BoardCast 将消息广播给 Message
func (server *Server) BoardCast(user *User, message string) {
	msg := "[" + user.Addr + "]: " + message
	server.Message <- msg
}

// Handler 用户上线后的主流程
func (server *Server) Handler(conn net.Conn) {
	//fmt.Printf("New connection from %s %s\n", conn.RemoteAddr(), conn.LocalAddr())
	user := NewUser(conn)

	server.mapLock.Lock()
	server.UserMap[user.Addr] = user
	server.mapLock.Unlock()

	server.BoardCast(user, "已上线")

	// 写操作
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := user.conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("写操作错误：", err)
				return
			}
			if n == 0 {
				server.BoardCast(user, "已下线")
				return
			}
			msg := buf[:n-1]
			server.BoardCast(user, string(msg))
		}

	}()

	//select {}
}
