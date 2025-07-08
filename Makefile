.PHONY: build run test clean fmt vet

# 默认目标
all: fmt vet test build

# 构建应用程序
build:
	@echo "构建应用程序..."
	@go build -o bin/server cmd/server/main.go

# 运行应用程序
run:
	@echo "启动服务器..."
	@go run cmd/server/main.go

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...

# 静态分析
vet:
	@echo "静态分析..."
	@go vet ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf bin/

# 创建bin目录
bin:
	@mkdir -p bin

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod tidy
	@go mod download

# 生成文档
doc:
	@echo "生成文档..."
	@godoc -http=:6060

# 交叉编译
build-linux: bin
	@echo "为Linux构建..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go

build-windows: bin
	@echo "为Windows构建..."
	@GOOS=windows GOARCH=amd64 go build -o bin/server-windows.exe cmd/server/main.go

build-all: build-linux build-windows build