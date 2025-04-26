# Gin 模板项目 (极简版)

这是一个基于 Gin 框架的 API 模板项目极简版，提供了用户管理和认证功能。

## 简化概述

本项目是原 Gin 模板的极简版本，只保留了用户管理相关功能，适合快速启动新项目或学习 Gin 框架。

## 核心功能

- **用户认证**: 注册和登录功能，JWT 认证
- **用户管理**: 用户的增删改查和密码修改

## 极简项目结构

```
gin-template/
├── controllers/     # 控制器：处理HTTP请求
│   ├── auth_controller.go    # 认证控制器
│   └── user_controller.go    # 用户控制器
├── middlewares/     # 中间件：横切关注点
│   ├── auth_middleware.go    # 认证中间件
│   ├── cors_middleware.go    # CORS中间件
│   ├── logger.go             # 日志中间件
│   ├── error_handler.go      # 错误处理中间件
│   └── error_recovery.go     # 恢复中间件
├── models/          # 数据模型：结构定义
│   └── user.go      # 用户模型
├── services/        # 服务层：业务逻辑
│   └── user_service.go    # 用户服务
├── repositories/    # 数据仓储：数据库交互
│   └── user_repository.go    # 用户仓储
├── utils/           # 工具函数：通用功能
├── config/          # 配置管理
└── routes/          # 路由定义
    └── routes.go    # 极简的路由配置
```

## API 极简清单

### 认证 API

- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/register` - 用户注册

### 用户 API

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/:id` - 获取用户详情
- `PUT /api/v1/users/:id` - 更新用户信息 (需认证)
- `DELETE /api/v1/users/:id` - 删除用户 (需认证)
- `POST /api/v1/users/change-password` - 修改密码 (需认证)

## 保留的中间件

- **Logger**: 日志记录
- **Recovery**: 错误恢复
- **ErrorHandler**: 错误处理
- **CorsMiddleware**: 跨域支持
- **AuthMiddleware**: 身份验证

## 启动应用

```bash
# 下载依赖
go mod download

# 运行应用
go run main.go
```

默认应用会运行在 `http://localhost:9999`

## Swagger 文档

访问 `http://localhost:9999/swagger/index.html` 查看 API 文档

## 此版本的优势

1. **聚焦核心**: 只保留用户认证和管理功能，易于理解
2. **低学习曲线**: 简单直观，适合初学者
3. **容易扩展**: 基于此模板可以逐步添加其他业务功能
4. **维护成本低**: 代码量少，易于维护和测试 