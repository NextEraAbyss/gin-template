# Gin API 项目模板

基于Go语言Gin框架的API项目模板，采用三层架构设计，提供安全、高性能的API开发框架。

## 核心特性

- **分层架构**：标准三层架构(Controllers, Services, Repositories)
- **安全性**：JWT认证、密码强度验证、请求速率限制、CORS配置
- **高性能**：Redis缓存、数据库连接池优化
- **开发友好**：统一日志、标准API响应格式、结构化错误处理
- **代码质量**：符合Go最佳实践、自动化测试支持

## 项目结构

```
.
├── config/           # 配置管理
├── controllers/      # 控制器 (表示层)
├── middlewares/      # 中间件
├── models/           # 数据模型
├── repositories/     # 数据访问层
├── routes/           # 路由设置
├── services/         # 业务逻辑层
├── utils/            # 工具函数
├── validation/       # 请求验证
├── internal/         # 内部包
│   ├── container/    # 依赖注入容器
│   ├── mysql/        # MySQL连接管理
│   └── redis/        # Redis缓存管理
├── main.go           # 应用入口
└── logs/             # 应用日志目录
```

## API响应格式

```json
{
    "code": 0,          // 状态码，0表示成功，非0表示错误
    "message": "成功",   // 状态消息
    "data": {}          // 响应数据
}
```

## 错误码设计

```
// 系统级状态码
CodeSuccess       = 0    // 成功
CodeInvalidParams = 1001 // 无效的参数
CodeUnauthorized  = 1002 // 未授权
CodeForbidden     = 1003 // 禁止访问
CodeNotFound      = 1004 // 资源不存在
CodeInternalError = 1005 // 内部错误

// 用户相关错误
CodeUserNotFound  = 2001 // 用户不存在
CodeUserExists    = 2002 // 用户已存在
CodePasswordError = 2003 // 密码错误
CodeTokenExpired  = 2004 // Token过期
```

## 快速开始

### 环境要求

- Go 1.16+
- MySQL 5.7+
- Redis

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/YourUsername/gin-template.git
```

2. 安装依赖
```bash
go mod download
```

3. 创建.env文件 (参考示例配置)

4. 运行项目
```bash
go run main.go
```

## 许可证

MIT