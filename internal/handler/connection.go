package handler

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"awesomeProject/pkg/models"
)

// ConnectionHandler 连接处理器
type ConnectionHandler struct {
	users map[string]*models.User
}

// NewConnectionHandler 创建新的连接处理器
func NewConnectionHandler() *ConnectionHandler {
	return &ConnectionHandler{
		users: make(map[string]*models.User),
	}
}

// Handle 处理单个连接
func (h *ConnectionHandler) Handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("关闭连接失败: %v", err)
		}
	}()

	log.Printf("新连接建立: %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	
	// 设置连接超时
	if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
		log.Printf("设置读取超时失败: %v", err)
		return
	}

	reader := bufio.NewReader(conn)
	
	// 发送欢迎消息
	if _, err := conn.Write([]byte("欢迎连接到服务器！请输入您的用户名: ")); err != nil {
		log.Printf("发送欢迎消息失败: %v", err)
		return
	}

	for {
		// 读取客户端消息
		message, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("客户端 %s 断开连接", conn.RemoteAddr())
			} else {
				log.Printf("读取消息失败: %v", err)
			}
			break
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		// 处理消息
		response := h.processMessage(message, conn.RemoteAddr().String())
		
		// 发送响应
		if _, err := conn.Write([]byte(response + "\n")); err != nil {
			log.Printf("发送响应失败: %v", err)
			break
		}

		// 重置读取超时
		if err := conn.SetReadDeadline(time.Now().Add(30 * time.Second)); err != nil {
			log.Printf("重置读取超时失败: %v", err)
			break
		}
	}
}

// processMessage 处理接收到的消息
func (h *ConnectionHandler) processMessage(message, clientAddr string) string {
	parts := strings.Fields(message)
	if len(parts) == 0 {
		return "无效的命令"
	}

	command := strings.ToLower(parts[0])
	
	switch command {
	case "register":
		if len(parts) < 3 {
			return "用法: register <用户名> <邮箱>"
		}
		return h.registerUser(parts[1], parts[2], clientAddr)
	
	case "login":
		if len(parts) < 2 {
			return "用法: login <用户名>"
		}
		return h.loginUser(parts[1], clientAddr)
	
	case "info":
		return h.getUserInfo(clientAddr)
	
	case "help":
		return h.getHelp()
		
	case "quit", "exit":
		return "再见！"
	
	default:
		return fmt.Sprintf("未知命令: %s。输入 'help' 查看帮助", command)
	}
}

// registerUser 注册用户
func (h *ConnectionHandler) registerUser(name, email, clientAddr string) string {
	if _, exists := h.users[clientAddr]; exists {
		return "您已经注册过了"
	}

	user := models.NewUser(name, email)
	if err := user.Validate(); err != nil {
		return fmt.Sprintf("注册失败: %v", err)
	}

	h.users[clientAddr] = user
	return fmt.Sprintf("注册成功！欢迎 %s", name)
}

// loginUser 用户登录
func (h *ConnectionHandler) loginUser(name, clientAddr string) string {
	if user, exists := h.users[clientAddr]; exists {
		if user.Name == name {
			return fmt.Sprintf("欢迎回来，%s！", name)
		}
		return "用户名不匹配"
	}
	return "请先注册账户"
}

// getUserInfo 获取用户信息
func (h *ConnectionHandler) getUserInfo(clientAddr string) string {
	if user, exists := h.users[clientAddr]; exists {
		return fmt.Sprintf("用户信息: %s", user.String())
	}
	return "请先注册或登录"
}

// getHelp 获取帮助信息
func (h *ConnectionHandler) getHelp() string {
	return `可用命令:
register <用户名> <邮箱> - 注册新用户
login <用户名> - 登录
info - 查看用户信息
help - 显示帮助
quit/exit - 退出`
}