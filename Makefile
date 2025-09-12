# LLM Demo Makefile

.PHONY: build run test clean help

# 默认目标
.DEFAULT_GOAL := help

# 构建项目
build:
	@echo "构建 LLM Demo 项目..."
	go build -o bin/llm-demo main.go

# 运行项目
run:
	@echo "启动 LLM Demo 服务..."
	go run main.go

# 运行测试
test:
	@echo "运行 API 测试..."
	@if [ ! -f bin/llm-demo ]; then \
		echo "请先构建项目: make build"; \
		exit 1; \
	fi
	@echo "启动服务进行测试..."
	@./bin/llm-demo &
	@SERVER_PID=$$!; \
	sleep 3; \
	./test_api.sh; \
	kill $$SERVER_PID

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/

# 安装依赖
deps:
	@echo "安装项目依赖..."
	go mod tidy
	go mod download

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
lint:
	@echo "运行代码检查..."
	golangci-lint run

# 生成 API 代码（如果使用 goctl）
gen-api:
	@echo "生成 API 代码..."
	goctl api go -api api/chat.api -dir . --style=goZero

# 显示帮助信息
help:
	@echo "LLM Demo 项目命令:"
	@echo "  build     - 构建项目"
	@echo "  run       - 运行项目"
	@echo "  test      - 运行 API 测试"
	@echo "  clean     - 清理构建文件"
	@echo "  deps      - 安装依赖"
	@echo "  fmt       - 格式化代码"
	@echo "  lint      - 代码检查"
	@echo "  gen-api   - 生成 API 代码"
	@echo "  help      - 显示此帮助信息"
