# TCP 服务器项目

这是一个使用 Go 语言开发的模块化 TCP 服务器项目，支持用户注册、登录等基本功能。

## 项目结构

```
.
├── cmd/
│   └── server/          # 应用程序入口点
│       └── main.go
├── pkg/                 # 可被外部导入的公共库
│   ├── config/          # 配置管理
│   ├── models/          # 数据模型
│   └── server/          # 服务器核心逻辑
├── internal/            # 内部模块
│   └── handler/         # 连接处理器
├── bin/                 # 编译后的二进制文件
├── Makefile            # 构建脚本
├── go.mod              # Go 模块文件
└── README.md           # 项目说明
```

## 功能特性

- **模块化设计**: 采用标准 Go 项目布局，代码职责分离清晰
- **配置管理**: 支持环境变量配置，灵活部署
- **优雅关闭**: 支持信号处理和优雅关闭
- **错误处理**: 完善的错误处理和日志记录
- **用户管理**: 支持用户注册、登录、信息查询
- **连接超时**: 自动处理连接超时和清理
- **并发安全**: 支持多客户端并发连接

## 快速开始

### 环境要求

- Go 1.21 或更高版本

### 安装依赖

```bash
make deps
```

### 构建项目

```bash
make build
```

### 运行服务器

```bash
make run
```

或者直接运行二进制文件：

```bash
./bin/server
```

### 运行测试

```bash
make test
```

## 配置

可以通过环境变量配置服务器参数：

```bash
export SERVER_HOST=0.0.0.0    # 服务器监听地址，默认 127.0.0.1
export SERVER_PORT=9999       # 服务器监听端口，默认 8888
```

## 使用示例

### 连接到服务器

使用 telnet 或 nc 命令连接：

```bash
telnet localhost 8888
# 或
nc localhost 8888
```

### 可用命令

连接后可以使用以下命令：

- `register <用户名> <邮箱>` - 注册新用户
- `login <用户名>` - 登录
- `info` - 查看用户信息
- `help` - 显示帮助
- `quit` 或 `exit` - 退出连接

### 示例会话

```
$ telnet localhost 8888
欢迎连接到服务器！请输入您的用户名: 
register zhangsan zhangsan@example.com
注册成功！欢迎 zhangsan

info
用户信息: User{ID: 0, Name: zhangsan, Email: zhangsan@example.com}

help
可用命令:
register <用户名> <邮箱> - 注册新用户
login <用户名> - 登录
info - 查看用户信息
help - 显示帮助
quit/exit - 退出

quit
再见！
```

## 开发

### 代码格式化

```bash
make fmt
```

### 静态分析

```bash
make vet
```

### 清理构建文件

```bash
make clean
```

### 交叉编译

```bash
make build-linux    # Linux 平台
make build-windows  # Windows 平台
make build-all      # 所有平台
```

## 架构说明

### 关键组件

1. **Server**: 核心服务器组件，负责监听连接和生命周期管理
2. **ConnectionHandler**: 连接处理器，负责处理客户端消息和业务逻辑
3. **Config**: 配置管理，支持环境变量和默认值
4. **Models**: 数据模型，定义用户等实体

### 设计原则

- **单一职责**: 每个模块都有明确的职责
- **依赖注入**: 通过构造函数注入依赖
- **接口隔离**: 使用接口定义组件间的交互
- **错误处理**: 统一的错误处理策略
- **可测试性**: 代码结构便于单元测试

## 扩展建议

1. **数据持久化**: 添加数据库支持（MySQL、PostgreSQL等）
2. **身份验证**: 实现 JWT 或会话管理
3. **消息队列**: 添加 Redis 或 RabbitMQ 支持
4. **监控指标**: 集成 Prometheus 监控
5. **配置热重载**: 支持配置文件热重载
6. **HTTP API**: 添加 REST API 支持
7. **WebSocket**: 支持 WebSocket 协议

## 许可证

MIT License
