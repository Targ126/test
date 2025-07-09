package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string

	Addr string

	C chan string

	conn net.Conn
}

func (u *User) Listen() {
	for {
		msg := <-u.C
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("用户写 error: ", err)
			continue
		}
	}
}

func NewUser(conn net.Conn) *User {
	user := &User{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		conn: conn,
	}
	go user.Listen()
	return user
}
