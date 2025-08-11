# Makefile for Feishu MCP Go

.PHONY: help build run clean test lint install deps

# 默认目标
help:
	@echo "可用命令:"
	@echo "  build    - 构建二进制文件"
	@echo "  run      - 运行程序"
	@echo "  clean    - 清理构建文件"
	@echo "  test     - 运行测试"
	@echo "  lint     - 代码检查"
	@echo "  install  - 安装到系统"
	@echo "  deps     - 下载依赖"

# 变量定义
APP_NAME = feishu-mcp-go
BUILD_DIR = build
CMD_DIR = cmd/server
VERSION ?= $(shell git describe --tags --always --dirty || echo "dev")
BUILD_TIME = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

# 构建
build: deps
	@echo "构建 $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)

# 运行
run: build
	@echo "运行 $(APP_NAME)..."
	./$(BUILD_DIR)/$(APP_NAME)

# 清理
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	go clean

# 测试
test:
	@echo "运行测试..."
	go test -v ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	@which golangci-lint > /dev/null || (echo "请安装 golangci-lint" && exit 1)
	golangci-lint run

# 安装到系统
install: build
	@echo "安装 $(APP_NAME) 到 /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/

# 跨平台构建
build-all: deps
	@echo "跨平台构建..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./$(CMD_DIR)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./$(CMD_DIR)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./$(CMD_DIR)

# 开发模式运行
dev: deps
	@echo "开发模式运行..."
	go run ./$(CMD_DIR) --log-level=debug

# 显示版本信息
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)" 