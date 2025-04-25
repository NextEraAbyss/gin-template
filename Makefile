.PHONY: all build run test clean lint help build-prod compress size-compare fmt check-deps security-check coverage bench race docs

# 设置变量
BINARY_NAME=gin-template
BINARY_UNIX=$(BINARY_NAME)_unix

# 默认目标
all: lint test build

# 构建应用
build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) main.go

# 生产环境构建
build-prod:
	@echo "Building for production..."
	powershell -Command "$$env:CGO_ENABLED=0; $$env:GOOS='linux'; $$env:GOARCH='amd64'; go build -trimpath -ldflags='-w -s -buildid=' -o bin/$(BINARY_NAME) main.go"
	@echo "Build complete. Binary location: bin/$(BINARY_NAME)"

# 压缩二进制文件
compress:
	@echo "Compressing binary with UPX..."
	upx --best --lzma bin/$(BINARY_NAME)
	@echo "Compression complete."

# 构建并压缩
build-compress: build-prod compress

# 显示压缩前后的大小比较
size-compare:
	@echo "Binary size comparison:"
	powershell -Command "ls bin/$(BINARY_NAME) | Select-Object Length,Name | Format-Table"

# 运行应用
run:
	@echo "Running..."
	go run main.go

# 运行测试
test:
	@echo "Running tests..."
	go test -v ./...

# 运行测试并生成覆盖率报告
coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 运行基准测试
bench:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# 运行竞态检测
race:
	@echo "Running race detector..."
	go test -race ./...

# 清理构建文件
clean:
	@echo "Cleaning..."
	powershell -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }"
	powershell -Command "if (Test-Path coverage.out) { Remove-Item coverage.out }"
	powershell -Command "if (Test-Path coverage.html) { Remove-Item coverage.html }"
	go clean

# 运行代码检查
lint:
	@echo "Running linters..."
	golangci-lint run

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# 检查依赖
check-deps:
	@echo "Checking dependencies..."
	go mod tidy
	go mod verify

# 安全检查
security-check:
	@echo "Running security checks..."
	gosec -quiet -exclude-dir=vendor -exclude-dir=docs -exclude-dir=bin ./...
	golangci-lint run --enable=gosec

# 安装依赖
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest

# 生成 API 文档
docs:
	@echo "Generating API documentation..."
	swag init --parseDependency --parseInternal
	@echo "API documentation generated successfully."

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all            - Run lint, test and build"
	@echo "  build          - Build the application"
	@echo "  build-prod     - Build for production with optimizations"
	@echo "  compress       - Compress binary using UPX"
	@echo "  build-compress - Build and compress for production"
	@echo "  size-compare   - Compare binary sizes"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage report"
	@echo "  bench          - Run benchmarks"
	@echo "  race           - Run race detector"
	@echo "  clean          - Clean build files"
	@echo "  lint           - Run code linters"
	@echo "  fmt            - Format code"
	@echo "  check-deps     - Check dependencies"
	@echo "  security-check - Run security checks"
	@echo "  deps           - Install dependencies"
	@echo "  docs           - Generate API documentation"
	@echo "  help           - Show this help message" 