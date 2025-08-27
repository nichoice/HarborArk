package controller

import (
	"HarborArk/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MonitorController 监控控制器
type MonitorController struct {
	monitorService *service.MonitorService
}

// NewMonitorController 创建监控控制器
func NewMonitorController(monitorService *service.MonitorService) *MonitorController {
	return &MonitorController{
		monitorService: monitorService,
	}
}

// GetSystemOverview 获取系统概览
// @Summary 获取系统监控概览
// @Description 获取CPU、内存、磁盘、网络的概览信息
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=service.SystemOverview}
// @Failure 500 {object} Response
// @Router /monitor/system/overview [get]
func (mc *MonitorController) GetSystemOverview(c *gin.Context) {
	overview, err := mc.monitorService.GetSystemOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    overview,
	})
}

// GetSystemMetrics 获取系统指标
// @Summary 获取系统指标数据
// @Description 获取指定时间范围内的系统指标数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param range query string false "时间范围" default(1h)
// @Success 200 {object} Response{data=monitor.SystemMetrics}
// @Failure 500 {object} Response
// @Router /monitor/system/metrics [get]
func (mc *MonitorController) GetSystemMetrics(c *gin.Context) {
	timeRange := c.DefaultQuery("range", "1h")

	metrics, err := mc.monitorService.GetSystemMetrics(timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    metrics,
	})
}

// GetMetricsByType 按类型获取指标
// @Summary 按类型获取指标数据
// @Description 获取指定类型和时间范围的指标数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param type path string true "指标类型" Enums(cpu,memory,disk,network)
// @Param range query string false "时间范围" default(1h)
// @Param limit query int false "限制数量" default(100)
// @Success 200 {object} Response{data=[]monitor.MetricRecord}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /monitor/metrics/{type} [get]
func (mc *MonitorController) GetMetricsByType(c *gin.Context) {
	metricType := c.Param("type")
	timeRange := c.DefaultQuery("range", "1h")
	limitStr := c.DefaultQuery("limit", "100")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	// 验证指标类型
	validTypes := map[string]bool{
		"cpu": true, "memory": true, "disk": true, "network": true,
		"nas": true, "protocol": true,
	}
	if !validTypes[metricType] {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的指标类型",
		})
		return
	}

	records, err := mc.monitorService.GetMetricsByType(metricType, timeRange)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	// 限制返回数量
	if len(records) > limit {
		records = records[:limit]
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    records,
	})
}

// GetLatestMetrics 获取最新指标
// @Summary 获取最新指标数据
// @Description 获取指定类型的最新指标数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param type query string true "指标类型"
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} Response{data=[]monitor.MetricRecord}
// @Failure 500 {object} Response
// @Router /monitor/latest [get]
func (mc *MonitorController) GetLatestMetrics(c *gin.Context) {
	metricType := c.Query("type")
	limitStr := c.DefaultQuery("limit", "10")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	records, err := mc.monitorService.GetLatestMetrics(metricType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    records,
	})
}

// GetMetricsByTimeRange 按时间范围获取指标
// @Summary 按时间范围获取指标
// @Description 获取指定时间范围内的指标数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param type query string true "指标类型"
// @Param name query string false "指标名称"
// @Param start query string true "开始时间" format(date-time)
// @Param end query string true "结束时间" format(date-time)
// @Success 200 {object} Response{data=[]monitor.MetricRecord}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /monitor/range [get]
func (mc *MonitorController) GetMetricsByTimeRange(c *gin.Context) {
	metricType := c.Query("type")
	metricName := c.Query("name")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if metricType == "" || startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "缺少必要参数",
		})
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的开始时间格式",
		})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的结束时间格式",
		})
		return
	}

	records, err := mc.monitorService.GetMetricsByTimeRange(metricType, metricName, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    records,
	})
}

// GetAggregatedMetrics 获取聚合指标
// @Summary 获取聚合指标数据
// @Description 获取指定时间范围内的聚合指标数据
// @Tags 监控
// @Accept json
// @Produce json
// @Param type query string true "指标类型"
// @Param interval query string false "聚合间隔" default(1h)
// @Param start query string true "开始时间" format(date-time)
// @Param end query string true "结束时间" format(date-time)
// @Success 200 {object} Response{data=[]storage.AggregatedMetric}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /monitor/aggregated [get]
func (mc *MonitorController) GetAggregatedMetrics(c *gin.Context) {
	metricType := c.Query("type")
	intervalStr := c.DefaultQuery("interval", "1h")
	startStr := c.Query("start")
	endStr := c.Query("end")

	if metricType == "" || startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "缺少必要参数",
		})
		return
	}

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的聚合间隔格式",
		})
		return
	}

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的开始时间格式",
		})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的结束时间格式",
		})
		return
	}

	records, err := mc.monitorService.GetAggregatedMetrics(metricType, interval, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    records,
	})
}

// GetMonitorHealth 获取监控服务健康状态
// @Summary 获取监控服务健康状态
// @Description 获取监控服务的运行状态和配置信息
// @Tags 监控
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=map[string]interface{}}
// @Router /monitor/health [get]
func (mc *MonitorController) GetMonitorHealth(c *gin.Context) {
	health := mc.monitorService.GetHealth()

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "ok",
		Data:    health,
	})
}
