# HarborArk

一个基于 Go 和 Gin 框架构建的现代化 RESTful API 项目，集成了完整的 Swagger 文档、高性能日志系统和灵活的配置管理。

## ✨ 特性

- 🚀 **高性能**: 基于 Gin 框架，提供高性能的 HTTP 服务
- 📚 **API 文档**: 集成 Swagger UI，自动生成交互式 API 文档
- 📝 **日志系统**: 基于 Zap 的结构化日志，支持文件轮转和多级别输出
- ⚙️ **配置管理**: 支持 YAML 配置文件和环境变量，配置热重载
- 🏗️ **模块化设计**: 清晰的项目结构，易于维护和扩展
- 🛡️ **中间件支持**: 内置日志记录、错误恢复等中间件

## 📁 项目结构

```
HarborArk/
├── cmd/                    # 命令行工具
│   ├── docs/              # Swagger 文档生成文件
│   ├── root.go            # Cobra 根命令
│   ├── server.go          # 服务器启动命令
│   └── version.go         # 版本信息
├── config/                # 配置管理
│   ├── config.go          # 基础配置
│   ├── setting.go         # 配置结构定义
│   └── settings-dev.yaml  # 开发环境配置
├── internal/              # 内部包
│   ├── controller/        # 控制器层
│   ├── model/            # 数据模型
│   ├── repository/       # 数据访问层
│   ├── service/          # 业务逻辑层
│   └── utils/            # 工具函数
├── router/               # 路由配置
│   ├── api/              # API 路由
│   ├── middleware/       # 中间件
│   └── swagger.go        # Swagger 路由配置
├── logs/                 # 日志文件
├── go.mod               # Go 模块文件
├── go.sum               # 依赖校验文件
├── main.go              # 程序入口
└── README.md            # 项目说明
```

## 🚀 快速开始

### 环境要求

- Go 1.25.0 或更高版本
- Git

### 安装依赖

```bash
# 克隆项目
git clone <your-repo-url>
cd HarborArk

# 安装依赖
go mod tidy

# 安装 Swagger 工具（用于生成文档）
go install github.com/swaggo/swag/cmd/swag@latest
```

### 运行项目

```bash
# 启动开发服务器
go run main.go server

# 或者构建后运行
go build -o harborark
./harborark server
```

服务器启动后，您可以访问：

- **主页**: http://localhost:8080/
- **API 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health
- **API 信息**: http://localhost:8080/api/info

## 📖 API 文档

项目集成了 Swagger UI，提供完整的 API 文档和在线测试功能。

### 访问文档

启动服务器后，访问 http://localhost:8080/swagger/index.html 查看完整的 API 文档。

### 使用 Cobra 命令管理文档

项目集成了 Cobra 命令行工具来管理 Swagger 文档：

```bash
# 查看 Swagger 命令帮助
go run main.go swagger --help

# 生成 Swagger 文档
go run main.go swagger generate

# 强制重新生成文档
go run main.go swagger generate --force

# 自定义输出目录和主文件
go run main.go swagger generate --output docs --main cmd/server.go

# 验证 Swagger 文档
go run main.go swagger validate

# 清理生成的文档
go run main.go swagger clean
```

### 自动更新文档

在配置文件中启用自动更新：

```yaml
swagger:
  autoUpdate: true  # 启用自动更新
```

启用后，服务器启动时会自动检查并更新 Swagger 文档。

### 示例 API

项目提供了用户管理的示例 API：

- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/{id}` - 根据 ID 获取用户
- `POST /api/v1/users` - 创建新用户

## ⚙️ 配置

### 配置文件

项目使用 YAML 格式的配置文件，位于 `config/settings-dev.yaml`：

```yaml
server:
  port: "8080"
  mode: "debug"

logger:
  level: debug
  encoding: json
  filename: logs/app.log
  maxSize: 100
  maxAge: 7
  maxBackups: 10

swagger:
  title: "HarborArk API"
  description: "HarborArk 项目 API 文档"
  version: "1.0.0"
  host: "localhost:8080"
  basePath: "/api/v1"
  enabled: true
```

### 环境变量

支持通过环境变量覆盖配置，环境变量前缀为 `HARBORARK_`：

```bash
export HARBORARK_SERVER_PORT=9000
export HARBORARK_LOGGER_LEVEL=info
```

## 📝 日志系统

项目使用 Zap 高性能日志库，支持：

- **结构化日志**: JSON 格式输出
- **多级别日志**: Debug、Info、Warn、Error、Fatal
- **文件轮转**: 自动按大小和时间轮转日志文件
- **双输出**: 开发模式下同时输出到控制台和文件

### 日志配置

```yaml
logger:
  level: debug          # 日志级别
  encoding: json        # 编码格式 (json/console)
  filename: logs/app.log # 日志文件路径
  maxSize: 100          # 单个文件最大大小 (MB)
  maxAge: 7             # 保留天数
  maxBackups: 10        # 保留文件数量
```

## 🛠️ 开发指南

### 添加新的 API

1. 在 `internal/controller/` 中创建控制器
2. 添加 Swagger 注释
3. 在 `cmd/server.go` 中注册路由
4. 重新生成文档

示例控制器：

```go
// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取所有用户的列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func GetUsers(c *gin.Context) {
    // 实现逻辑
}
```

### 中间件

项目内置了以下中间件：

- **GinLogger**: HTTP 请求日志记录
- **GinRecovery**: Panic 恢复和错误处理

添加自定义中间件：

```go
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 中间件逻辑
        c.Next()
    }
}
```

## 🧪 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/controller

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📦 构建和部署

### 本地构建

```bash
# 构建二进制文件
go build -o harborark

# 交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o harborark-linux

# 交叉编译 (Windows)
GOOS=windows GOARCH=amd64 go build -o harborark.exe
```

### Docker 部署

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o harborark

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/harborark .
COPY --from=builder /app/config ./config
CMD ["./harborark", "server"]
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🔗 相关链接

- [Gin 框架文档](https://gin-gonic.com/)
- [Swagger 文档](https://swagger.io/)
- [Zap 日志库](https://github.com/uber-go/zap)
- [Viper 配置库](https://github.com/spf13/viper)
- [Cobra CLI 库](https://github.com/spf13/cobra)

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue: [GitHub Issues](https://github.com/your-username/HarborArk/issues)
- 邮箱: your-email@example.com

---

⭐ 如果这个项目对您有帮助，请给它一个 Star！