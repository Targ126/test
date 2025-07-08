package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"awesomeProject/internal/handler"
	"awesomeProject/pkg/config"
)

// Server TCP服务器
type Server struct {
	config    *config.Config
	listener  net.Listener
	handler   *handler.ConnectionHandler
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// New 创建新的服务器实例
func New(cfg *config.Config) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Server{
		config:  cfg,
		handler: handler.NewConnectionHandler(),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	// 创建监听器
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("启动服务器失败: %w", err)
	}
	
	s.listener = listener
	log.Printf("服务器启动成功，监听地址: %s", addr)

	// 设置信号处理，支持优雅关闭
	s.setupSignalHandler()

	// 开始接受连接
	s.wg.Add(1)
	go s.acceptConnections()

	// 等待所有goroutine完成
	s.wg.Wait()
	log.Println("服务器已关闭")
	
	return nil
}

// Stop 停止服务器
func (s *Server) Stop() error {
	log.Println("正在关闭服务器...")
	
	// 取消上下文
	s.cancel()
	
	// 关闭监听器
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("关闭监听器失败: %w", err)
		}
	}
	
	return nil
}

// acceptConnections 接受新连接
func (s *Server) acceptConnections() {
	defer s.wg.Done()
	
	for {
		select {
		case <-s.ctx.Done():
			log.Println("停止接受新连接")
			return
		default:
			// 设置Accept超时，以便能够检查context
			if tcpListener, ok := s.listener.(*net.TCPListener); ok {
				tcpListener.SetDeadline(time.Now().Add(1 * time.Second))
			}
			
			conn, err := s.listener.Accept()
			if err != nil {
				// 检查是否是超时错误
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				
				// 如果context已取消，正常退出
				if s.ctx.Err() != nil {
					return
				}
				
				log.Printf("接受连接失败: %v", err)
				continue
			}
			
			// 处理连接
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection 处理单个连接
func (s *Server) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	
	// 在context取消时关闭连接
	go func() {
		<-s.ctx.Done()
		conn.Close()
	}()
	
	s.handler.Handle(conn)
}

// setupSignalHandler 设置信号处理器，支持优雅关闭
func (s *Server) setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		log.Println("收到关闭信号，正在优雅关闭服务器...")
		s.Stop()
	}()
}