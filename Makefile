.PHONY: all build run test clean lint help

# 默认目标
all: lint test build

# 构建应用
build:
	@echo "Building..."
	go build -o bin/app main.go

# 运行应用
run:
	@echo "Running..."
	go run main.go

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

# 运行代码检查
lint:
	@echo "Running linters..."
	golangci-lint run

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 生成 API 文档
docs:
	@echo "Generating API documentation..."
	swag init

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all        - Run lint, test and build"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  clean      - Clean build files"
	@echo "  lint       - Run code linters"
	@echo "  deps       - Install dependencies"
	@echo "  docs       - Generate API documentation"
	@echo "  migrate    - Run database migrations"
	@echo "  help       - Show this help message" 