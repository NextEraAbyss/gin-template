# Gin API 项目模板

这是一个使用Go语言的Gin框架开发的API项目模板，遵循标准三层架构设计。该模板提供了一套完整的、可扩展的API开发框架，具有良好的安全性、性能和可维护性。

## 项目特性

- **分层架构**：严格的三层架构设计，使代码结构清晰，职责分明
- **安全性**：
  - JWT认证机制
  - 密码强度验证
  - 可配置的令牌过期时间
  - 请求速率限制
  - CORS安全配置
- **性能优化**：
  - Redis缓存支持
  - 数据库连接池优化
  - 高效的查询处理
  - 异步日志处理
- **开发体验**：
  - 统一日志系统
  - 标准化API响应格式
  - 结构化错误处理
  - 请求数据验证
  - 完整的开发工具链
- **代码质量**：
  - 符合Go最佳实践
  - 代码注释完善
  - 模块化设计
  - 静态代码分析
  - 自动化测试支持

## 项目架构

本项目采用经典的三层架构设计：

1. **表示层（Controllers）**：处理HTTP请求和响应
2. **业务逻辑层（Services）**：实现业务逻辑
3. **数据访问层（Repositories）**：处理数据存取

### 数据流动示意图

```
HTTP请求
   ↓
Controllers (表示层)
   ↓  
Services (业务逻辑层)
   ↓
Repositories (数据访问层)
   ↓
数据库 (MySQL/Redis)
```

## 项目结构

```
.
├── config/           # 配置相关代码
│   └── config.go     # 配置管理
├── controllers/      # 控制器/处理器 (表示层)
│   ├── auth_controller.go     # 认证相关控制器
│   └── user_controller.go     # 用户相关控制器
├── middlewares/      # 中间件
│   ├── auth_middleware.go     # JWT认证中间件
│   ├── cors_middleware.go     # 跨域请求处理
│   ├── error_handler.go       # 全局错误处理
│   ├── logger.go              # 请求日志记录
│   └── rate_limit.go          # 速率限制中间件
├── models/           # 数据模型
│   ├── auth.go                # 认证相关模型
│   └── user.go                # 用户模型
├── repositories/     # 数据访问层
│   └── user_repository.go     # 用户数据访问
├── routes/           # 路由设置
│   └── routes.go              # API路由配置
├── services/         # 业务逻辑层
│   └── user_service.go        # 用户业务逻辑
├── utils/            # 工具函数
│   ├── jwt.go                 # JWT工具
│   ├── logger.go              # 日志工具
│   ├── password.go            # 密码处理
│   └── response.go            # 统一响应格式
├── validation/       # 请求验证
│   └── user_dto.go            # 用户相关DTO
├── internal/         # 内部包
│   ├── container/     # 依赖注入容器
│   │   └── container.go       # 依赖注入管理
│   ├── mysql/        # MySQL连接管理
│   │   └── mysql.go           # 数据库连接与管理
│   └── redis/        # Redis缓存管理
│       └── redis.go           # Redis连接与操作
├── main.go           # 应用入口
└── logs/             # 应用日志目录(运行时创建)
```

## 编码规范

### 目录和文件命名

- 目录名使用**小写单词**
- Go源文件使用**小写单词**，如果包含多个单词，使用下划线分隔
- 接口和struct命名使用**大驼峰**

### 包结构

- **models**: 定义数据库模型
- **validation**: 定义DTO(数据传输对象)和验证规则
- **controllers**: 处理HTTP请求，参数校验，调用业务逻辑
- **services**: 实现业务逻辑
- **repositories**: 实现数据访问
- **utils**: 工具类
- **middlewares**: 中间件

### 接口设计

所有的服务和仓库都应该定义接口，以便于测试和解耦：

```go
// BaseService 所有服务的基础接口
type BaseService interface {
    ServiceName() string
}

// UserService 用户服务接口
type UserService interface {
    BaseService
    // 具体方法...
}
```

### 错误处理

- 使用自定义错误码和错误消息进行错误处理
- 错误返回格式统一
- 敏感信息不直接返回给客户端

### 数据传输对象(DTO)

- 所有的请求和响应都使用DTO对象
- DTO定义在`validation`包中
- 模型和DTO分离，避免暴露敏感字段

## API响应格式

统一的API响应格式：

```json
{
    "code": 0,       // 状态码，0表示成功，非0表示错误
    "message": "操作成功", // 状态消息
    "data": {}       // 响应数据
}
```

## 认证机制

本项目使用JWT(JSON Web Token)进行身份验证：

1. 用户登录成功后，服务器生成JWT
2. 客户端在后续请求的Authorization头中携带Token
3. 服务器通过中间件验证Token的有效性

## 错误码设计

```
// 系统级状态码 (1000-1999)
CodeSuccess       = 0    // 成功
CodeInvalidParams = 1001 // 无效的参数
CodeUnauthorized  = 1002 // 未授权
CodeForbidden     = 1003 // 禁止访问
CodeNotFound      = 1004 // 资源不存在
CodeInternalError = 1005 // 内部错误
CodeServerError   = 1006 // 服务器错误

// 用户相关错误 (2000-2999)
CodeUserNotFound  = 2001 // 用户不存在
CodeUserExists    = 2002 // 用户已存在
CodePasswordError = 2003 // 密码错误
CodeTokenExpired  = 2004 // Token过期
CodeTokenInvalid  = 2005 // Token无效
CodeUserDisabled  = 2006 // 用户已禁用
```

## 快速开始

### 环境要求

- Go 1.16+
- MySQL 5.7+
- Redis (可选)

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/YourUsername/gin-template.git
```

2. 安装依赖
```bash
go mod download
```

3. 创建并配置.env文件
```
# 服务器配置
ENV=development
SERVER_HOST=0.0.0.0
SERVER_PORT=9999

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=gin_template

# Redis配置
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT配置
JWT_SECRET=your_jwt_secret
JWT_EXPIRATION_HOURS=24
```

4. 运行项目
```bash
go run main.go
```

### 部署

1. 构建项目
```bash
go build -o gin-template main.go
```

2. 配置.env文件
```
ENV=production
...其他配置
```

3. 运行
```bash
./gin-template
```

## 测试

```bash
go test ./... -v
```

## 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

MIT