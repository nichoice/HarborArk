package service

import (
	"HarborArk/config"
	"HarborArk/internal/monitor"
	"HarborArk/internal/monitor/collector"
	"HarborArk/internal/monitor/storage"
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// MonitorService 监控服务
type MonitorService struct {
	monitor *monitor.Monitor
	storage *storage.DBStorage
	config  config.MonitorConfig
}

// NewMonitorService 创建监控服务
func NewMonitorService(db *gorm.DB, cfg config.MonitorConfig) *MonitorService {
	// 创建存储
	dbStorage := storage.NewDBStorage(db)

	// 创建监控实例
	monitorInstance := monitor.NewMonitor(cfg)
	monitorInstance.SetStorage(dbStorage)
	monitorInstance.SetScheduler(monitor.NewDefaultScheduler())

	// 添加采集器
	systemCollector := collector.NewSystemCollector()
	monitorInstance.AddCollector(systemCollector)

	return &MonitorService{
		monitor: monitorInstance,
		storage: dbStorage,
		config:  cfg,
	}
}

// Start 启动监控服务
func (s *MonitorService) Start() error {
	if !s.config.Enabled {
		zap.L().Info("监控服务已禁用")
		return nil
	}

	return s.monitor.Start()
}

// Stop 停止监控服务
func (s *MonitorService) Stop() error {
	return s.monitor.Stop()
}

// IsRunning 检查是否运行中
func (s *MonitorService) IsRunning() bool {
	return s.monitor.IsRunning()
}

// GetSystemMetrics 获取系统指标
func (s *MonitorService) GetSystemMetrics(timeRange string) (*monitor.SystemMetrics, error) {
	_, err := time.ParseDuration(timeRange)
	if err != nil {
		// 使用默认5分钟，这里暂时不使用duration
	}

	return s.monitor.GetSystemOverview()
}

// GetMetricsByType 按类型获取指标
func (s *MonitorService) GetMetricsByType(metricType string, timeRange string) ([]monitor.MetricRecord, error) {
	duration, err := time.ParseDuration(timeRange)
	if err != nil {
		duration = 1 * time.Hour // 默认1小时
	}

	return s.monitor.GetMetrics(monitor.MetricType(metricType), duration)
}

// GetLatestMetrics 获取最新指标
func (s *MonitorService) GetLatestMetrics(metricType string, limit int) ([]monitor.MetricRecord, error) {
	ctx := context.Background()
	return s.storage.GetLatestMetrics(ctx, metricType, limit)
}

// GetMetricsByTimeRange 按时间范围获取指标
func (s *MonitorService) GetMetricsByTimeRange(metricType, metricName string, start, end time.Time) ([]monitor.MetricRecord, error) {
	ctx := context.Background()
	return s.storage.GetMetricsByTimeRange(ctx, metricType, metricName, start, end)
}

// GetAggregatedMetrics 获取聚合指标
func (s *MonitorService) GetAggregatedMetrics(metricType string, interval time.Duration, start, end time.Time) ([]storage.AggregatedMetric, error) {
	ctx := context.Background()
	return s.storage.AggregateMetrics(ctx, metricType, interval, start, end)
}

// GetSystemOverview 获取系统概览
func (s *MonitorService) GetSystemOverview() (*SystemOverview, error) {
	ctx := context.Background()

	// 获取最新的各类指标
	cpuRecords, _ := s.storage.GetLatestMetrics(ctx, string(monitor.MetricTypeCPU), 10)
	memoryRecords, _ := s.storage.GetLatestMetrics(ctx, string(monitor.MetricTypeMemory), 10)
	diskRecords, _ := s.storage.GetLatestMetrics(ctx, string(monitor.MetricTypeDisk), 10)
	networkRecords, _ := s.storage.GetLatestMetrics(ctx, string(monitor.MetricTypeNetwork), 10)

	overview := &SystemOverview{
		CPU:       s.buildCPUOverview(cpuRecords),
		Memory:    s.buildMemoryOverview(memoryRecords),
		Disk:      s.buildDiskOverview(diskRecords),
		Network:   s.buildNetworkOverview(networkRecords),
		Timestamp: time.Now(),
	}

	return overview, nil
}

// SystemOverview 系统概览
type SystemOverview struct {
	CPU       CPUOverview     `json:"cpu"`
	Memory    MemoryOverview  `json:"memory"`
	Disk      DiskOverview    `json:"disk"`
	Network   NetworkOverview `json:"network"`
	Timestamp time.Time       `json:"timestamp"`
}

// CPUOverview CPU概览
type CPUOverview struct {
	Usage    float64 `json:"usage"`
	LoadAvg1 float64 `json:"load_avg_1"`
	LoadAvg5 float64 `json:"load_avg_5"`
	Cores    int     `json:"cores"`
}

// MemoryOverview 内存概览
type MemoryOverview struct {
	Usage     float64 `json:"usage"`
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
}

// DiskOverview 磁盘概览
type DiskOverview struct {
	Usage float64 `json:"usage"`
	Total uint64  `json:"total"`
	Used  uint64  `json:"used"`
}

// NetworkOverview 网络概览
type NetworkOverview struct {
	BytesRecv uint64 `json:"bytes_recv"`
	BytesSent uint64 `json:"bytes_sent"`
}

// buildCPUOverview 构建CPU概览
func (s *MonitorService) buildCPUOverview(records []monitor.MetricRecord) CPUOverview {
	overview := CPUOverview{}

	for _, record := range records {
		switch record.MetricName {
		case "cpu_usage_percent":
			overview.Usage = record.Value
		case "load_avg_1":
			overview.LoadAvg1 = record.Value
		case "load_avg_5":
			overview.LoadAvg5 = record.Value
		case "cpu_cores":
			overview.Cores = int(record.Value)
		}
	}

	return overview
}

// buildMemoryOverview 构建内存概览
func (s *MonitorService) buildMemoryOverview(records []monitor.MetricRecord) MemoryOverview {
	overview := MemoryOverview{}

	for _, record := range records {
		switch record.MetricName {
		case "memory_usage_percent":
			overview.Usage = record.Value
		case "memory_total_bytes":
			overview.Total = uint64(record.Value)
		case "memory_used_bytes":
			overview.Used = uint64(record.Value)
		case "memory_available_bytes":
			overview.Available = uint64(record.Value)
		}
	}

	return overview
}

// buildDiskOverview 构建磁盘概览
func (s *MonitorService) buildDiskOverview(records []monitor.MetricRecord) DiskOverview {
	overview := DiskOverview{}

	// 简化处理：取第一个磁盘的数据
	for _, record := range records {
		switch record.MetricName {
		case "disk_usage_percent":
			if overview.Usage == 0 { // 只取第一个
				overview.Usage = record.Value
			}
		case "disk_total_bytes":
			if overview.Total == 0 {
				overview.Total = uint64(record.Value)
			}
		case "disk_used_bytes":
			if overview.Used == 0 {
				overview.Used = uint64(record.Value)
			}
		}
	}

	return overview
}

// buildNetworkOverview 构建网络概览
func (s *MonitorService) buildNetworkOverview(records []monitor.MetricRecord) NetworkOverview {
	overview := NetworkOverview{}

	// 聚合所有网络接口的数据
	for _, record := range records {
		switch record.MetricName {
		case "network_receive_bytes":
			overview.BytesRecv += uint64(record.Value)
		case "network_transmit_bytes":
			overview.BytesSent += uint64(record.Value)
		}
	}

	return overview
}

// Cleanup 清理过期数据
func (s *MonitorService) Cleanup() error {
	ctx := context.Background()
	return s.storage.Cleanup(ctx)
}

// GetHealth 获取监控服务健康状态
func (s *MonitorService) GetHealth() map[string]interface{} {
	return map[string]interface{}{
		"running":   s.IsRunning(),
		"enabled":   s.config.Enabled,
		"interval":  s.config.CollectInterval,
		"timestamp": time.Now(),
	}
}
