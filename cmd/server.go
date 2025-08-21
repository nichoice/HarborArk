package cmd

import (
	"HarborArk/config"
	"HarborArk/internal/i18n"
	"HarborArk/internal/middleware"
	"HarborArk/internal/migration"
	"HarborArk/router"
	api "HarborArk/router/api"
	routerMiddleware "HarborArk/router/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	_ "HarborArk/cmd/docs"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "启动 HarborArk 服务器",
	Long:  `启动 HarborArk API 服务器`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	serverCmd.Flags().StringP("port", "p", "8080", "服务器运行端口")
	rootCmd.AddCommand(serverCmd)
}

// @title HarborArk API
// @version 1.0
// @description HarborArk 项目 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
func startServer() {
	// 初始化配置
	if err := config.Init(); err != nil {
		panic(fmt.Errorf("初始化配置失败: %v", err))
	}

	// 获取配置
	serverConfig := config.GetServerConfig()
	logConfig := config.GetLogConfig()
	swaggerConfig := config.GetSwaggerConfig()

	// 初始化日志系统
	if err := routerMiddleware.Init(logConfig, serverConfig.Mode); err != nil {
		panic(fmt.Errorf("初始化日志失败: %v", err))
	}

	// 初始化国际化
	if err := i18n.Init("translations"); err != nil {
		zap.L().Fatal("国际化初始化失败", zap.Error(err))
	}
	zap.L().Info("国际化初始化成功")

	// 初始化数据库
	config.ConnectDB()

	// 数据库迁移
	if err := migration.AutoMigrate(); err != nil {
		zap.L().Fatal("数据库迁移失败", zap.Error(err))
	}

	// 自动更新 Swagger 文档
	if swaggerConfig.AutoUpdate && swaggerConfig.Enabled {
		AutoUpdateSwaggerDocs()
	}

	// 设置 Gin 模式
	gin.SetMode(serverConfig.Mode)
	if serverConfig.Mode == "debug" {
		gin.ForceConsoleColor()
	}

	// 创建路由
	r := gin.New()

	// 添加中间件
	r.Use(middleware.CORSMiddleware()) // 添加跨域中间件
	r.Use(routerMiddleware.GinLogger(), routerMiddleware.GinRecovery(true))

	// 设置 Swagger 文档
	router.SetupSwagger(r)

	// 设置用户管理路由
	api.SetupUserRoutes(r)

	// 基础路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": i18n.T(c, "welcome"),
			"version": "1.0.0",
			"docs":    "/swagger/index.html",
		})
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": i18n.T(c, "health_ok"),
			"time":    "2025-01-21T11:00:00Z",
		})
	})

	// 启动服务器
	port := ":" + serverConfig.Port
	zap.L().Info("服务器启动中...",
		zap.String("port", serverConfig.Port),
		zap.String("mode", serverConfig.Mode),
		zap.String("docs", "http://localhost:"+serverConfig.Port+"/swagger/index.html"),
	)

	if err := r.Run(port); err != nil {
		zap.L().Fatal("服务器启动失败", zap.Error(err))
	}
}
