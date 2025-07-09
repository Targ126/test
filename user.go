package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name string

	Addr string

	C chan string

	conn net.Conn

	server *Server
	
	// 添加done channel用于优雅关闭
	done chan struct{}
}

func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		C:      make(chan string),
		conn:   conn,
		server: server,
		done:   make(chan struct{}),
	}
	go user.Listen()
	return user
}

func (user *User) sendMyself(msg string) {
	select {
	case user.C <- msg:
	case <-user.done:
		// 如果用户已经断开，直接返回
		return
	}
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		user.server.mapLock.Lock()
		for _, onlineUser := range user.server.UserMap {
			user.sendMyself(onlineUser.Name)
		}
		user.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := user.server.UserMap[newName]
		if ok {
			user.sendMyself(newName + "已经被使用，请重试")
			return
		}

		user.server.mapLock.Lock()
		delete(user.server.UserMap, user.Name)
		user.Name = newName
		user.server.UserMap[newName] = user
		user.server.mapLock.Unlock()

		user.sendMyself("已成功将名字修改为：" + newName)
	} else {
		user.server.BoardCase(user, msg)
	}

	//user.server.Message <- fmt.Sprintf("[" + user.Addr + "] " + user.Name + " : " + msg)
}

func (user *User) online() {

	user.server.mapLock.Lock()
	user.server.UserMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.DoMessage("已上线")
}

func (user *User) offline() {
	user.server.mapLock.Lock()
	delete(user.server.UserMap, user.Name)
	user.server.mapLock.Unlock()

	user.DoMessage("已下线")
}

// 添加优雅关闭函数
func (user *User) Close() {
	// 先从用户表中移除
	user.server.mapLock.Lock()
	delete(user.server.UserMap, user.Name)
	user.server.mapLock.Unlock()
	
	// 关闭done channel，通知Listen函数停止
	close(user.done)
	
	// 关闭连接
	user.conn.Close()
	
	// 最后关闭消息channel
	close(user.C)
}

func (u *User) Listen() {
	defer func() {
		// 确保连接被关闭
		u.conn.Close()
	}()
	
	for {
		select {
		case msg, ok := <-u.C:
			if !ok {
				// channel已关闭，退出
				return
			}
			_, err := u.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("接收广播错误: ", err)
				return
			}
		case <-u.done:
			// 收到关闭信号，退出
			return
		}
	}
}
