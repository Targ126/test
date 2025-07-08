package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (server *Server) Start() {
	//监听
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	//关闭
	defer listener.Close()
	//
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		server.Handler(conn)
	}
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Printf("New connection from %s %s\n", conn.RemoteAddr(), conn.LocalAddr())
}
