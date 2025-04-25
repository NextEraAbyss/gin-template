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

## 项目结构

```
.
├── config/         # 配置相关代码
├── controllers/    # 控制器/处理器 (表示层)
├── middleware/     # 中间件
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

### Windows 环境配置

1. 安装 Chocolatey (Windows 包管理器)
```powershell
# 以管理员身份运行 PowerShell，执行以下命令
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

2. 安装必要的工具
```powershell
# 安装 Make
choco install make

# 安装 UPX (用于压缩二进制文件)
choco install upx
```

3. 验证安装
```powershell
# 检查 Make 版本
make --version

# 检查 UPX 版本
upx --version
```

### 安装

1. 克隆项目
```bash
git clone https://github.com/NextEraAbyss/gin-template.git
cd gin-template
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置环境变量
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑.env文件，设置你的环境变量
```

4. 运行项目
```bash
# 开发模式运行
make run

# 或者直接运行
go run main.go
```

### 开发命令

项目提供了 Makefile 来简化开发流程：

```bash
# 构建应用（开发环境）
make build

# 构建应用（生产环境，包含优化）
make build-prod

# 构建并压缩二进制文件
make build-compress

# 比较二进制文件大小
make size-compare

# 运行应用
make run

# 运行测试
make test

# 运行测试并生成覆盖率报告（生成 coverage.html）
make coverage

# 运行基准测试
make bench

# 运行竞态检测
make race

# 运行代码检查
make lint

# 格式化代码
make fmt

# 检查依赖
make check-deps

# 运行安全检查
make security-check

# 安装依赖和工具
make deps

# 生成 API 文档
make docs

# 清理构建文件
make clean

# 显示帮助信息
make help
```

每个命令的具体功能：

| 命令 | 说明 |
|------|------|
| `build` | 构建开发环境应用 |
| `build-prod` | 构建生产环境应用（包含优化参数） |
| `build-compress` | 构建并使用 UPX 压缩二进制文件 |
| `size-compare` | 显示压缩前后的二进制文件大小比较 |
| `run` | 运行应用 |
| `test` | 运行单元测试 |
| `coverage` | 运行测试并生成 HTML 格式的覆盖率报告 |
| `bench` | 运行基准测试，测试代码性能 |
| `race` | 运行竞态检测，发现并发问题 |
| `lint` | 运行代码静态检查工具 |
| `fmt` | 格式化代码，保持代码风格一致 |
| `check-deps` | 检查并验证项目依赖 |
| `security-check` | 运行安全检查，发现潜在安全问题 |
| `deps` | 安装项目依赖和开发工具 |
| `docs` | 生成 Swagger API 文档 |
| `clean` | 清理构建产物和临时文件 |
| `help` | 显示所有可用的 make 命令 |

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

### 业务状态码

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

### 用户相关

**获取用户列表**
```
GET /api/users
```

查询参数:
- `page`: 页码，默认1
- `page_size`: 每页条数，默认10
- `sort`: 排序字段
- `order`: 排序方向(asc或desc)
- `search`: 搜索关键词

**获取单个用户**
```
GET /api/users/:id
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

**更新用户** (需要认证)
```
PUT /api/users/:id
```

**删除用户** (需要认证)
```
DELETE /api/users/:id
```

### 文章相关

**获取文章列表**
```
GET /api/articles
```

查询参数:
- `page`: 页码，默认1
- `page_size`: 每页条数，默认10
- `keyword`: 搜索关键词
- `status`: 文章状态(1:草稿,2:已发布)
- `order_by`: 排序字段(created_at,view_count)
- `order`: 排序方向(asc,desc)

**获取单个文章**
```
GET /api/articles/:id
```

**创建文章** (需要认证)
```
POST /api/articles
```

请求体:
```json
{
  "title": "新文章标题",
  "content": "新文章内容",
  "status": 1
}
```

**更新文章** (需要认证)
```
PUT /api/articles/:id
```

**删除文章** (需要认证)
```
DELETE /api/articles/:id
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

## 部署指南

### 开发环境部署

1. 安装依赖
```bash
make deps
```

2. 配置环境变量
```bash
cp .env.example .env
# 编辑 .env 文件
```

3. 运行项目
```bash
make run
```

### 生产环境部署

#### Docker部署

1. 构建镜像
```bash
docker build -t gin-api .
```

2. 运行容器
```bash
docker run -d -p 8080:8080 --name gin-api gin-api
```

#### 传统部署

1. 构建应用
```bash
make build
```

2. 使用PM2运行（推荐）
```bash
pm2 start ecosystem.config.js
```

3. 使用nohup运行
```bash
./start.sh
```

### Nginx配置

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

## 开发规范

### 代码风格

- 使用 `golangci-lint` 进行代码检查
- 遵循 Go 官方代码规范
- 使用 `go fmt` 格式化代码
- 保持一致的命名风格

### Git提交规范

提交信息格式：
```
<type>(<scope>): <subject>

<body>

<footer>
```

类型（type）：
- feat: 新功能
- fix: 修复bug
- docs: 文档更新
- style: 代码格式调整
- refactor: 重构
- test: 测试相关
- chore: 构建过程或辅助工具的变动

### 测试规范

- 单元测试覆盖率要求 > 80%
- 集成测试覆盖主要业务流程
- 使用 `go test -cover` 检查测试覆盖率

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

MIT License

## 联系方式

- 项目维护者：[秦若宸/NextEraAbyss]
- 邮箱：[1578347363@qq.com]
- 项目链接：[https://github.com/NextEraAbyss/gin-template](https://github.com/NextEraAbyss/gin-template)