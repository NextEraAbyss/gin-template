# 🚀 Gin 用户API极简模板

![Go Version](https://img.shields.io/badge/Go-1.23.0+-00ADD8?style=flat-square&logo=go)
![Gin Version](https://img.shields.io/badge/Gin-1.9+-00ADD8?style=flat-square&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-5.7+-4479A1?style=flat-square&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-3.0+-DC382D?style=flat-square&logo=redis&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=flat-square)

基于Go语言Gin框架的极简API项目模板，仅保留用户相关功能，采用三层架构设计，提供安全、高性能的API开发基础。

## ✨ 核心特性 | Features

### 🏗️ 架构设计
- **🔄 分层架构** - 标准三层架构(Controllers, Services, Repositories)
- **💉 依赖注入** - 优雅的依赖管理和解耦
- **🔧 查询构建器** - 链式调用，支持复杂查询构建
- **📦 模块化设计** - 清晰的项目结构，易于维护扩展

### 🔒 安全特性
- **🛡️ JWT认证** - 无状态身份认证，安全可靠
- **🚦 速率限制** - 防止API滥用（默认每分钟1000次）
- **🌐 CORS配置** - 灵活的跨域资源共享设置
- **🔐 安全HTTP头** - 自动添加安全相关响应头
- **🔑 密码加密** - bcrypt加密存储，支持强度验证

### ⚡ 性能优化
- **🗄️ Redis缓存** - 自动缓存查询结果，提升响应速度
- **🔗 连接池优化** - 数据库连接池管理，提高并发性能
- **📊 慢查询监控** - 自动检测和记录性能问题
- **💾 批量操作** - 支持批量创建、更新、删除操作

### 🛠️ 开发友好
- **📝 统一日志** - 结构化日志记录，支持文件输出
- **📋 标准响应** - 统一的API响应格式
- **🎯 错误处理** - 结构化错误处理和状态码管理
- **📚 Swagger文档** - 自动生成API接口文档
- **🔍 代码质量** - golangci-lint检查，符合Go最佳实践

### 🚀 现代化工具
- **⚙️ Make脚本** - 便捷的构建和开发命令
- **🐳 容器化支持** - Docker部署配置
- **🔄 优雅关闭** - 支持服务优雅停止
- **📈 性能监控** - 请求追踪和性能统计

## 📊 性能指标 | Performance

- ⚡ **启动时间** < 3秒
- 🎯 **API响应** < 100ms (本地环境)
- 📈 **并发处理** 1000+ QPS
- 💾 **内存占用** < 50MB (基础运行)
- 🔄 **缓存命中率** > 80%

## 🏗️ 项目结构 | Project Structure

```
.
├── 📁 controllers/              # 🎮 控制器层 - HTTP请求处理
│   └── user_controller.go       #   用户控制器
├── 📁 services/                 # 💼 业务逻辑层 - 核心业务处理  
│   └── user_service.go          #   用户业务逻辑
├── 📁 repositories/             # 🗄️ 数据访问层 - 数据库操作
│   └── user_repository.go       #   用户数据访问
├── 📁 models/                   # 📋 数据模型 - 实体定义
│   └── user.go                  #   用户实体模型
├── 📁 routes/                   # 🛤️ 路由配置 - API路径定义
│   └── routes.go                #   路由设置
├── 📁 middlewares/              # 🔧 中间件 - 通用功能处理
│   ├── auth_middleware.go       #   JWT认证中间件
│   ├── cors_middleware.go       #   CORS跨域中间件
│   ├── logger.go                #   日志记录中间件
│   ├── rate_limit.go            #   速率限制中间件
│   └── security.go              #   安全头中间件
├── 📁 internal/                 # 🔒 内部模块
│   ├── container/               #   📦 依赖注入容器
│   ├── mysql/                   #   🗄️ 数据库连接和查询构建器
│   └── redis/                   #   🔴 Redis连接管理
├── 📁 config/                   # ⚙️ 配置管理
├── 📁 utils/                    # 🛠️ 工具函数
├── 📁 validation/               # ✅ 数据验证
├── 📁 docs/                     # 📚 Swagger文档
├── 📄 main.go                   # 🚀 应用入口
├── 📄 Makefile                  # 🔨 构建脚本
├── 📄 .golangci.yml             # 🔍 代码质量配置
├── 📄 go.mod                    # 📦 Go模块定义
└── 📄 README.md                 # 📖 项目说明
```

## 🏛️ 架构设计 | Architecture

### 🔄 三层架构模式
- **🎮 Controller Layer** - 处理HTTP请求，参数验证，响应格式化
- **💼 Service Layer** - 业务逻辑处理，事务管理，缓存策略
- **🗄️ Repository Layer** - 数据访问抽象，SQL查询，数据转换

### 🧩 核心组件

#### 🔗 查询构建器 (QueryBuilder)
- **⛓️ 链式调用** - 支持复杂查询的链式构建
- **💾 缓存集成** - 自动缓存查询结果
- **📊 性能监控** - 慢查询检测和统计
- **🔗 连接池管理** - 数据库连接池优化

#### 🔧 中间件链
- **🔍 请求ID追踪** - 全链路请求跟踪
- **📝 日志记录** - 结构化日志输出
- **🛡️ 错误恢复** - 优雅的错误处理
- **📋 统一错误处理** - 标准化错误响应
- **🌐 CORS跨域** - 跨域资源共享配置
- **🔒 安全HTTP头** - 安全相关HTTP头设置
- **🚦 速率限制** - 防止API滥用（默认每分钟1000次）

## 🛤️ 用户API路由 | API Routes

| 方法   | 路径                              | 说明         | 认证 | 描述                    |
|--------|-----------------------------------|--------------|------|-------------------------|
| 📄 GET    | `/api/v1/users`                   | 获取用户列表 | ❌   | 支持分页、搜索、排序    |
| 👤 GET    | `/api/v1/users/:id`               | 获取单个用户 | ❌   | 根据ID获取用户详情      |
| ✏️ PUT    | `/api/v1/users/:id`               | 更新用户信息 | ✅   | 更新用户资料信息        |
| 🗑️ DELETE | `/api/v1/users/:id`               | 删除用户     | ✅   | 软删除用户记录          |
| 🔑 POST   | `/api/v1/users/change-password`   | 修改密码     | ✅   | 修改当前用户密码        |

> 🔑 需要认证的接口需携带有效JWT Token：`Authorization: Bearer <token>`

### 🔍 查询参数 (GET /api/v1/users)
- `page` - 📖 页码，从1开始 (默认: 1)
- `page_size` - 📊 每页记录数 (默认: 10)
- `keyword` - 🔍 搜索关键词，支持用户名、邮箱模糊搜索
- `order_by` - 📈 排序字段 (默认: id)
- `order` - 🔀 排序方向 asc/desc (默认: desc)

## 📝 API响应格式 | Response Format

### 标准响应格式
```json
{
    "code": 0,              // 状态码，0表示成功，非0表示错误
    "message": "成功",       // 状态消息
    "data": {}              // 响应数据
}
```

### 📊 分页响应格式
```json
{
    "code": 0,
    "message": "成功", 
    "data": {
        "total": 100,       // 总记录数
        "items": [...],     // 当前页数据
        "page": 1,          // 当前页码
        "page_size": 10,    // 每页大小
        "pages": 10         // 总页数
    }
}
```

## 🚨 错误码设计 | Error Codes

```go
// 🏷️ 系统级状态码
CodeSuccess       = 0    // ✅ 成功
CodeInvalidParams = 1001 // ❌ 无效的参数
CodeUnauthorized  = 1002 // 🚫 未授权
CodeForbidden     = 1003 // 🛑 禁止访问
CodeNotFound      = 1004 // 🔍 资源不存在
CodeInternalError = 1005 // ⚠️ 内部错误

// 👥 用户相关错误
CodeUserNotFound  = 2001 // 👤 用户不存在
CodeUserExists    = 2002 // 👥 用户已存在
CodePasswordError = 2003 // 🔐 密码错误
CodeTokenExpired  = 2004 // ⏰ Token过期
```

## 📚 Swagger文档 | API Documentation

🌐 访问 **[http://localhost:9999/swagger/index.html](http://localhost:9999/swagger/index.html)** 查看完整的API接口文档。

## 🚀 快速开始 | Quick Start

### 📋 环境要求 | Requirements

- ![Go](https://img.shields.io/badge/Go-1.23.0+-00ADD8?style=flat-square&logo=go) Go 1.23.0+
- ![MySQL](https://img.shields.io/badge/MySQL-5.7+-4479A1?style=flat-square&logo=mysql&logoColor=white) MySQL 5.7+
- ![Redis](https://img.shields.io/badge/Redis-3.0+-DC382D?style=flat-square&logo=redis&logoColor=white) Redis 3.0+
- 🔨 Make 工具（Windows用户需要单独安装）

### 🪟 Windows 环境配置 | Windows Setup

#### 🔨 安装 Make 工具

Windows 系统默认不包含 Make 工具，推荐以下安装方式：

**通过 Chocolatey 安装（推荐）**

1. 安装 Chocolatey（以管理员身份运行 PowerShell）：
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

2. 安装 Make：
```powershell
choco install make
```

#### ✅ 验证安装

安装完成后，重新打开命令行窗口，验证安装：
```bash
make --version
```

### 💻 本地开发 | Local Development

#### 1. **📥 克隆项目**
```bash
git clone https://github.com/YourUsername/gin-template.git
cd gin-template
```

#### 2. **📦 安装依赖**
```bash
go mod download
```

#### 3. **⚙️ 配置环境变量**

创建 `.env` 文件（参考配置）：
```env
# 🌍 环境配置
ENV=development

# 🖥️ 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=9999

# 🗄️ 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gin_template

# 🔴 Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# 🔑 JWT配置
JWT_SECRET=your-super-secret-key-change-this-in-production
JWT_EXPIRATION_HOURS=24
```

#### 4. **🗄️ 数据库初始化**
```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE gin_template;"

# 运行项目时会自动迁移表结构（开发环境）
```

#### 5. **🔨 使用 Make 命令运行项目**
```bash
# 🚀 运行项目
make run

# 📦 构建项目
make build

# 🧪 运行测试
make test

# 🔍 代码检查
make lint

# 💄 代码格式化
make fmt

# 📚 生成API文档
make docs

# 📋 查看所有可用命令
make help
```

#### 6. **🐹 或者直接使用 Go 命令**
```bash
go run main.go
```

### 🛠️ 开发工具 | Development Tools

#### 🔨 Make 命令说明
- `make run` - 🚀 启动开发服务器
- `make build` - 📦 构建可执行文件
- `make build-prod` - 🏭 生产环境构建（优化）
- `make test` - 🧪 运行单元测试
- `make coverage` - 📊 生成测试覆盖率报告
- `make lint` - 🔍 代码质量检查
- `make fmt` - 💄 代码格式化
- `make docs` - 📚 生成Swagger文档
- `make clean` - 🧹 清理构建文件

#### 🔍 代码质量
项目使用 `golangci-lint` 进行代码质量检查，包括：
- ✨ 代码格式检查
- 🐛 潜在错误检测  
- ⚡ 性能问题识别
- 🔒 安全漏洞扫描
- 📊 代码复杂度分析

## ⚙️ 配置管理 | Configuration

### 📋 支持的配置源
- 🌍 环境变量
- 📄 `.env` 文件
- 🔧 默认值

### 🎛️ 主要配置项
- **🖥️ 服务器配置** - 端口、主机、超时设置
- **🗄️ 数据库配置** - 连接池、超时、字符集
- **🔴 Redis配置** - 连接池、超时、重试
- **🔑 JWT配置** - 密钥、过期时间
- **📝 日志配置** - 级别、输出格式

## 🚀 性能特性 | Performance Features

### 🗄️ 数据库优化
- **🔗 连接池管理** - 自动管理数据库连接
- **💾 查询缓存** - Redis缓存频繁查询
- **🐌 慢查询监控** - 自动检测和记录慢查询
- **📦 批量操作** - 支持批量创建、更新、删除

### 🚀 缓存策略
- **🔍 查询缓存** - 自动缓存数据库查询结果
- **🔄 缓存失效** - 数据更新时自动清理相关缓存
- **🔑 缓存键管理** - 统一的缓存键命名规范

### 📊 监控和日志
- **📝 结构化日志** - 使用Zap进行高性能日志记录
- **🔍 请求追踪** - 每个请求分配唯一ID
- **⏱️ 性能监控** - HTTP请求耗时统计
- **🚨 错误监控** - 自动捕获和记录错误

## 🔒 安全特性 | Security Features

### 🔐 认证和授权
- **🎫 JWT Token** - 无状态身份认证
- **🔐 密码加密** - 使用bcrypt加密存储
- **⏰ Token过期** - 自动处理Token过期

### 🛡️ 安全防护
- **🌐 CORS配置** - 跨域请求控制
- **🔒 安全头** - 自动添加安全相关HTTP头
- **🚦 速率限制** - 防止API滥用
- **✅ 输入验证** - 严格的参数验证

### 🔐 数据保护
- **💉 SQL注入防护** - 使用参数化查询
- **🔗 XSS防护** - 输入输出过滤
- **📁 文件权限** - 合理的文件访问权限

## 🌍 浏览器兼容性 | Browser Compatibility

### ✅ 支持的浏览器

- **Chrome** 60+ ✅
- **Firefox** 55+ ✅
- **Safari** 12+ ✅
- **Edge** 79+ ✅
- **移动浏览器** ✅

## 🚀 部署建议 | Deployment

### 🏭 生产环境配置
```bash
# 📦 构建生产版本
make build-prod

# ⚙️ 设置环境变量
export ENV=production
export JWT_SECRET=your-production-secret

# 🚀 启动服务
./bin/gin-template
```

### 🐳 Docker 部署
```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN make build-prod

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/gin-template .
CMD ["./gin-template"]
```

### 🌐 Nginx配置示例
```nginx
server {
    listen 80;
    server_name api.yourdomain.com;
    
    location / {
        proxy_pass http://localhost:9999;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## 🤝 贡献指南 | Contributing

欢迎提交Issue和Pull Request！

### 🎯 贡献方式

1. 🐛 **Bug报告** - 发现问题请提交Issue
2. 💡 **功能建议** - 欢迎提出改进建议
3. 🔧 **代码改进** - 提交Pull Request
4. 📖 **文档完善** - 改进项目文档

### 🔄 开发流程

```bash
# 1. 🍴 Fork项目
# 2. 🌿 创建特性分支
git checkout -b feature/your-feature

# 3. 💾 提交更改
git commit -m "Add: your feature"

# 4. 📤 推送分支
git push origin feature/your-feature

# 5. 🔄 创建Pull Request
```

### 📋 代码规范
- 📚 遵循 Go 官方代码规范
- 🔍 运行 `make lint` 确保代码质量
- 🧪 添加必要的单元测试
- 📖 更新相关文档

## 📄 许可证 | License

本项目采用 **MIT License** 开源协议 - 详见 [LICENSE](LICENSE) 文件

## 🔗 相关链接 | Links

- 📚 **API文档** - [http://localhost:9999/swagger/index.html](http://localhost:9999/swagger/index.html)
- 🐹 **Go官方文档** - [https://golang.org/doc/](https://golang.org/doc/)
- 🌟 **Gin框架** - [https://gin-gonic.com/](https://gin-gonic.com/)
- 🗄️ **GORM文档** - [https://gorm.io/](https://gorm.io/)
- 🔴 **Redis文档** - [https://redis.io/documentation](https://redis.io/documentation)

---
<div align="center">

**🌟 如果这个项目对您有帮助，请给个 Star 支持一下！**

<p>
  <img src="https://img.shields.io/github/stars/NextEraAbyss/cv?style=social" alt="GitHub stars">
  <img src="https://img.shields.io/github/forks/NextEraAbyss/cv?style=social" alt="GitHub forks">
</p>

*Built with ❤️ by [秦若宸](https://cv.wat.ink)*

**💼 极简而不简单，高效开发API项目的最佳选择！**

</div>