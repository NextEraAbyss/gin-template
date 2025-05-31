# Gin 用户API极简模板

基于Go语言Gin框架的极简API项目模板，仅保留用户相关功能，采用三层架构设计，提供安全、高性能的API开发基础。

## 核心特性

- **分层架构**：标准三层架构(Controllers, Services, Repositories)
- **安全性**：JWT认证、请求速率限制、CORS配置、安全HTTP头
- **高性能**：Redis缓存、数据库连接池优化
- **开发友好**：统一日志、标准API响应格式、结构化错误处理
- **代码质量**：符合Go最佳实践、自动化测试支持

## 项目结构（精简版）

```
.
├── controllers/      # 用户控制器
├── middlewares/      # 全局中间件
├── models/           # 用户数据模型
├── repositories/     # 用户数据访问层
├── routes/           # 路由设置（仅用户相关）
├── services/         # 用户业务逻辑层
├── internal/         # 依赖注入、数据库、缓存
├── main.go           # 应用入口
└── README.md         # 项目说明
```

## 全局中间件

- 请求ID追踪
- 日志记录
- 错误恢复
- 统一错误处理
- CORS跨域
- 安全HTTP头
- 速率限制（默认每分钟1000次）

## 用户API路由

| 方法   | 路径                        | 说明         | 认证 |
| ------ | --------------------------- | ------------ | ---- |
| GET    | /api/v1/users               | 获取用户列表 | 否   |
| GET    | /api/v1/users/:id           | 获取单个用户 | 否   |
| PUT    | /api/v1/users/:id           | 更新用户信息 | 是   |
| DELETE | /api/v1/users/:id           | 删除用户     | 是   |
| POST   | /api/v1/users/change-password | 修改密码   | 是   |

> 需要认证的接口需携带有效JWT Token。

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

## Swagger文档

访问 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) 查看API接口文档。

## 快速开始

### 环境要求

- Go 1.16+
- MySQL 5.7+
- Redis
- Make 工具（Windows用户需要单独安装）

### Windows 环境配置

#### 安装 Make 工具

Windows 系统默认不包含 Make 工具，推荐以下安装方式：

**方式一：通过 Chocolatey 安装（推荐）**

1. 安装 Chocolatey（以管理员身份运行 PowerShell）：
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

2. 安装 Make：
```powershell
choco install make
```

**方式二：通过 Scoop 安装**

1. 安装 Scoop：
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

2. 安装 Make：
```powershell
scoop install make
```

**方式三：下载预编译版本**

从 [GnuWin32](http://gnuwin32.sourceforge.net/packages/make.htm) 下载并安装，然后将安装路径添加到系统 PATH 环境变量中。

#### 验证安装

安装完成后，重新打开命令行窗口，验证安装：
```bash
make --version
```

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

4. 使用 Make 命令运行项目
```bash
# 运行项目
make run

# 构建项目
make build

# 运行测试
make test

# 代码检查
make lint

# 查看所有可用命令
make help
```

5. 或者直接使用 Go 命令
```bash
go run main.go
```

## 许可证

MIT