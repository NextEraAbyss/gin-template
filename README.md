# Gin API 项目模板

这是一个使用Go语言的Gin框架开发的API项目模板，遵循标准三层架构设计。该模板提供了一套完整的、可扩展的API开发框架，具有良好的安全性、性能和可维护性。

## 项目特性

- **分层架构**：严格的三层架构设计，使代码结构清晰，职责分明
- **安全性**：
  - JWT认证机制
  - 密码强度验证
  - 可配置的令牌过期时间
- **性能优化**：
  - Redis缓存支持
  - 数据库连接池优化
  - 高效的查询处理
- **开发体验**：
  - 统一日志系统
  - 标准化API响应格式
  - 结构化错误处理
  - 请求数据验证
- **代码质量**：
  - 符合Go最佳实践
  - 代码注释完善
  - 模块化设计

## 项目结构

```
.
├── config/         # 配置相关代码
├── controllers/    # 控制器/处理器 (表示层)
├── middlewares/    # 中间件
├── models/         # 数据模型和DTO
├── repositories/   # 数据访问层
├── routes/         # 路由设置
├── services/       # 业务逻辑层
├── utils/          # 工具函数
├── internal/       # 内部包
│   ├── database/   # 数据库管理
│   └── redis/      # Redis缓存管理
└── logs/           # 应用日志目录(运行时创建)
```

## 三层架构说明

1. **表示层**（Controllers）：处理HTTP请求，进行参数验证，并将任务委托给业务逻辑层。使用DTO(数据传输对象)进行请求验证和数据转换。
2. **业务逻辑层**（Services）：实现所有业务规则和逻辑，包括缓存处理、数据验证和业务操作。
3. **数据访问层**（Repositories）：处理数据持久化和检索，与数据库交互。

## 开始使用

### 前提条件

- Go 1.16+
- MySQL 5.7+
- Redis 6.0+ (可选，但推荐使用)

### 安装

1. 克隆项目
```bash
git clone https://github.com/your-username/gin-template.git
cd gin-template
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境变量
```bash
cp .env.example .env
# 编辑.env文件，设置你的环境变量
```

4. 运行项目
```bash
go run main.go
```

### 环境变量配置说明

项目使用`.env`文件或环境变量进行配置，主要配置项包括：

```
# 应用环境设置
ENV=development # 可选: development, production

# 服务器设置
SERVER_HOST=localhost
SERVER_PORT=8080

# 数据库设置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=gin_template

# Redis设置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT设置
JWT_SECRET=your-secure-jwt-secret-key-change-in-production
JWT_EXPIRATION_HOURS=24
```

## API 文档

### 标准响应格式

所有API响应遵循统一的JSON格式：

```json
{
  "code": 20000,       // 业务状态码
  "message": "操作成功", // 提示信息
  "data": {},          // 响应数据
  "error": ""          // 错误信息(仅在错误时返回)
}
```

### 业务状态码说明

| 状态码 | 描述 |
|-------|------|
| 20000 | 操作成功 |
| 20400 | 无内容 |
| 40001 | 请求参数错误 |
| 40100 | 未授权 |
| 40300 | 禁止访问 |
| 40400 | 资源不存在 |
| 50000 | 系统内部错误 |

### 认证相关

**登录**
```
POST /login
```

请求体:
```json
{
  "username": "your-username",
  "password": "your-password"
}
```

响应:
```json
{
  "code": 20000,
  "message": "操作成功",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "管理员",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### 用户相关

**获取所有用户**
```
GET /api/users
```

查询参数:
- `page`: 页码，默认1
- `page_size`: 每页条数，默认10
- `sort`: 排序字段
- `order`: 排序方向(asc或desc)
- `search`: 搜索关键词

响应:
```json
{
  "code": 20000,
  "message": "操作成功",
  "data": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "full_name": "管理员",
      "created_at": "2023-01-01T12:00:00Z",
      "updated_at": "2023-01-01T12:00:00Z"
    }
  ]
}
```

**获取单个用户**
```
GET /api/users/:id
```

响应:
```json
{
  "code": 20000,
  "message": "操作成功",
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "full_name": "管理员",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  }
}
```

**创建用户**
```
POST /api/users
```

请求体:
```json
{
  "username": "new-user",
  "email": "user@example.com",
  "password": "Password123!",
  "full_name": "New User"
}
```

注意：密码必须符合强度要求（至少8个字符，包含大小写字母、数字和特殊字符）

响应:
```json
{
  "code": 20000,
  "message": "操作成功",
  "data": {
    "id": 2,
    "username": "new-user",
    "email": "user@example.com",
    "full_name": "New User",
    "created_at": "2023-01-02T12:00:00Z",
    "updated_at": "2023-01-02T12:00:00Z"
  }
}
```

**更新用户** (需要认证)
```
PUT /api/users/:id
```

请求头:
```
Authorization: Bearer your-jwt-token
```

请求体:
```json
{
  "email": "updated-email@example.com",
  "full_name": "Updated Name"
}
```

响应:
```json
{
  "code": 20000,
  "message": "用户信息更新成功",
  "data": {
    "id": 2,
    "username": "new-user",
    "email": "updated-email@example.com",
    "full_name": "Updated Name",
    "created_at": "2023-01-02T12:00:00Z",
    "updated_at": "2023-01-02T12:10:00Z"
  }
}
```

**删除用户** (需要认证)
```
DELETE /api/users/:id
```

请求头:
```
Authorization: Bearer your-jwt-token
```

响应:
```json
{
  "code": 20000,
  "message": "用户删除成功",
  "data": null
}
```

## 日志系统

项目使用自定义日志系统，支持不同级别的日志记录和日志文件轮转：

- 日志级别：DEBUG, INFO, WARN, ERROR, FATAL
- 日志保存在 `./logs` 目录下
- 按日期生成不同的日志文件
- 错误日志单独存储
- 生产环境禁用DEBUG级别的日志

## 缓存系统

项目使用Redis进行数据缓存：

- 用户信息缓存
- 可配置的缓存过期时间
- 自动缓存失效处理
- 缓存命中率监控

## 打包和部署

### 在Windows上打包Linux版本

```bash
# 方法1: 直接使用环境变量（PowerShell）
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o app-linux-amd64 main.go

# 方法2: 直接使用环境变量（CMD）
set GOOS=linux
set GOARCH=amd64
go build -o app-linux-amd64 main.go

# 方法3: 使用go env命令（推荐）
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -o app-linux-amd64 main.go

# 完成后恢复为Windows环境（如果使用了方法3）
go env -w GOOS=windows
go env -w GOARCH=amd64

# 优化: 减小二进制文件大小
# 在上述任何一种方法中添加-ldflags参数
go build -ldflags="-s -w" -o app-linux-amd64 main.go
```


### 生产环境部署

#### 系统服务部署 (Linux)

1. 创建系统服务文件（/etc/systemd/system/gin-api.service）

```
[Unit]
Description=Gin API Template
After=network.target

[Service]
User=www-data
WorkingDirectory=/opt/gin-api
ExecStart=/opt/gin-api/app
Restart=on-failure
RestartSec=5s
Environment=ENV=production

[Install]
WantedBy=multi-user.target
```

2. 启用和管理服务

```bash
# 复制二进制文件和配置到部署目录
sudo mkdir -p /opt/gin-api
sudo cp app .env /opt/gin-api/
sudo chown -R www-data:www-data /opt/gin-api

# 启用服务
sudo systemctl enable gin-api

# 启动服务
sudo systemctl start gin-api

# 检查状态
sudo systemctl status gin-api

# 查看日志
sudo journalctl -u gin-api -f
```

#### Nginx反向代理配置

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 云平台部署

本项目也可以部署在各种云平台上：

- **AWS Elastic Beanstalk**
- **Google Cloud Run**
- **Azure App Service**
- **Heroku**

每个平台都有特定的部署流程，通常需要创建一个配置文件（如Procfile、app.yaml等）。详细部署步骤请参考相应云平台的官方文档。

## 贡献

欢迎贡献代码、提交问题和功能请求。请确保代码符合Go的最佳实践和项目的代码风格。

## 许可证

MIT License