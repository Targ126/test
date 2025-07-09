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
			// 使用select避免阻塞
			select {
			case user.C <- msg:
			case <-user.done:
				// 用户已断开，跳过
				continue
			default:
				// channel满了，跳过这个用户
				fmt.Printf("用户 %s 的消息队列已满，跳过发送\n", user.Name)
			}
		}
		server.mapLock.Unlock()
	}
}

// BoardCase 将消息广播给 Message
func (server *Server) BoardCase(user *User, message string) {
	server.Message <- fmt.Sprintf("[" + user.Addr + "] " + user.Name + " : " + message)
}

// Handler 用户上线后的主流程
func (server *Server) Handler(conn net.Conn) {
	//fmt.Printf("New connection from %s %s\n", conn.RemoteAddr(), conn.LocalAddr())
	user := NewUser(conn, server)

	user.online()

	isLive := make(chan bool)

	// 对写操作进行监听
	go func() {
		defer func() {
			// 确保在goroutine退出时清理用户
			user.offline()
		}()
		
		buf := make([]byte, 1024)
		for {
			n, err := user.conn.Read(buf)
			if err != nil && err != io.EOF {
				fmt.Println("用户读取错误：", err)
				return
			}
			if n == 0 {
				// 连接已关闭
				return
			}
			msg := buf[:n-1]
			user.DoMessage(string(msg))
			
			// 通知主循环用户还活着
			select {
			case isLive <- true:
			case <-user.done:
				// 用户已被踢，退出
				return
			}
		}
	}()

	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 5):
			user.sendMyself("长时间没动作，你被踢了")
			
			// 给用户一点时间接收最后的消息
			time.Sleep(100 * time.Millisecond)
			
			// 使用新的Close方法优雅关闭
			user.Close()
			return
		}
	}
	//select {}
}
