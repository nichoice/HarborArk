package api

import (
	"HarborArk/config"
	"HarborArk/internal/controller"
	"HarborArk/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupMonitorRoutes 设置监控路由
func SetupMonitorRoutes(r *gin.Engine, db *gorm.DB) {
	// 获取监控配置
	monitorConfig := config.GetMonitorConfig()

	// 创建监控服务
	monitorService := service.NewMonitorService(db, monitorConfig)

	// 启动监控服务
	if err := monitorService.Start(); err != nil {
		// 记录错误但不阻止服务启动
		// zap.L().Error("启动监控服务失败", zap.Error(err))
	}

	// 创建监控控制器
	monitorController := controller.NewMonitorController(monitorService)

	// 监控路由组
	monitorGroup := r.Group("/api/v1/monitor")
	{
		// 系统监控
		monitorGroup.GET("/system/overview", monitorController.GetSystemOverview)
		monitorGroup.GET("/system/metrics", monitorController.GetSystemMetrics)

		// 指标查询
		monitorGroup.GET("/metrics/:type", monitorController.GetMetricsByType)
		monitorGroup.GET("/latest", monitorController.GetLatestMetrics)
		monitorGroup.GET("/range", monitorController.GetMetricsByTimeRange)
		monitorGroup.GET("/aggregated", monitorController.GetAggregatedMetrics)

		// 服务状态
		monitorGroup.GET("/health", monitorController.GetMonitorHealth)
	}
}
