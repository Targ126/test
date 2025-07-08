package main

import (
	"log"

	"awesomeProject/pkg/config"
	"awesomeProject/pkg/server"
)

func main() {
	// 加载配置
	cfg := config.Load()
	
	// 创建服务器
	srv := server.New(cfg)
	
	// 启动服务器
	if err := srv.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}