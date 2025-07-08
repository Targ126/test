package config

import (
	"os"
	"strconv"
)

// Config 应用程序配置
type Config struct {
	Server ServerConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port int
}

// Load 从环境变量或默认值加载配置
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "127.0.0.1"),
			Port: getEnvAsInt("SERVER_PORT", 8888),
		},
	}
	return config
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量作为整数，如果不存在则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}