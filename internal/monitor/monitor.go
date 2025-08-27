package monitor

import (
	"HarborArk/config"
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Monitor 监控系统主类
type Monitor struct {
	config     config.MonitorConfig
	collectors []Collector
	storage    Storage
	scheduler  Scheduler
	running    bool
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewMonitor 创建新的监控实例
func NewMonitor(cfg config.MonitorConfig) *Monitor {
	ctx, cancel := context.WithCancel(context.Background())

	return &Monitor{
		config:     cfg,
		collectors: make([]Collector, 0),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// AddCollector 添加采集器
func (m *Monitor) AddCollector(collector Collector) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if collector.IsEnabled() {
		m.collectors = append(m.collectors, collector)
		zap.L().Info("添加采集器",
			zap.String("type", string(collector.GetType())))
	}
}

// SetStorage 设置存储
func (m *Monitor) SetStorage(storage Storage) {
	m.storage = storage
}

// SetScheduler 设置调度器
func (m *Monitor) SetScheduler(scheduler Scheduler) {
	m.scheduler = scheduler
}

// Start 启动监控系统
func (m *Monitor) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.config.Enabled {
		zap.L().Info("监控系统已禁用")
		return nil
	}

	if m.running {
		return fmt.Errorf("监控系统已在运行")
	}

	if m.storage == nil {
		return fmt.Errorf("存储未设置")
	}

	if m.scheduler == nil {
		return fmt.Errorf("调度器未设置")
	}

	if len(m.collectors) == 0 {
		return fmt.Errorf("没有可用的采集器")
	}

	// 解析采集间隔
	interval, err := time.ParseDuration(m.config.CollectInterval)
	if err != nil {
		return fmt.Errorf("无效的采集间隔: %v", err)
	}

	m.scheduler.SetInterval(interval)

	// 启动调度器
	if err := m.scheduler.Start(m.ctx, m.collect); err != nil {
		return fmt.Errorf("启动调度器失败: %v", err)
	}

	m.running = true

	zap.L().Info("监控系统启动成功",
		zap.String("interval", m.config.CollectInterval),
		zap.Int("collectors", len(m.collectors)))

	return nil
}

// Stop 停止监控系统
func (m *Monitor) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return nil
	}

	// 取消上下文
	m.cancel()

	// 停止调度器
	if m.scheduler != nil {
		if err := m.scheduler.Stop(); err != nil {
			zap.L().Error("停止调度器失败", zap.Error(err))
		}
	}

	m.running = false

	zap.L().Info("监控系统已停止")
	return nil
}

// IsRunning 检查是否运行中
func (m *Monitor) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// collect 执行数据采集
func (m *Monitor) collect() error {
	start := time.Now()

	var allRecords []MetricRecord
	var collectErrors []error

	// 并发采集所有指标
	recordsChan := make(chan []MetricRecord, len(m.collectors))
	errorsChan := make(chan error, len(m.collectors))

	var wg sync.WaitGroup

	for _, collector := range m.collectors {
		wg.Add(1)
		go func(c Collector) {
			defer wg.Done()

			records, err := c.Collect(m.ctx)
			if err != nil {
				errorsChan <- fmt.Errorf("采集器 %s 失败: %v", c.GetType(), err)
				return
			}

			recordsChan <- records
		}(collector)
	}

	// 等待所有采集器完成
	go func() {
		wg.Wait()
		close(recordsChan)
		close(errorsChan)
	}()

	// 收集结果
	for records := range recordsChan {
		allRecords = append(allRecords, records...)
	}

	// 收集错误
	for err := range errorsChan {
		collectErrors = append(collectErrors, err)
	}

	// 记录采集错误但不中断流程
	if len(collectErrors) > 0 {
		for _, err := range collectErrors {
			zap.L().Error("数据采集错误", zap.Error(err))
		}
	}

	// 存储数据
	if len(allRecords) > 0 {
		if err := m.storage.Store(m.ctx, allRecords); err != nil {
			zap.L().Error("存储数据失败", zap.Error(err))
			return err
		}
	}

	duration := time.Since(start)
	zap.L().Debug("数据采集完成",
		zap.Int("records", len(allRecords)),
		zap.Duration("duration", duration),
		zap.Int("errors", len(collectErrors)))

	return nil
}

// GetMetrics 获取最新指标
func (m *Monitor) GetMetrics(metricType MetricType, duration time.Duration) ([]MetricRecord, error) {
	if m.storage == nil {
		return nil, fmt.Errorf("存储未初始化")
	}

	endTime := time.Now()
	startTime := endTime.Add(-duration)

	req := QueryRequest{
		MetricType: string(metricType),
		StartTime:  startTime,
		EndTime:    endTime,
		Limit:      1000, // 默认限制
	}

	return m.storage.Query(m.ctx, req)
}

// GetSystemOverview 获取系统概览
func (m *Monitor) GetSystemOverview() (*SystemMetrics, error) {
	// 获取最近5分钟的数据
	duration := 5 * time.Minute

	cpuRecords, err := m.GetMetrics(MetricTypeCPU, duration)
	if err != nil {
		return nil, fmt.Errorf("获取CPU指标失败: %v", err)
	}

	memoryRecords, err := m.GetMetrics(MetricTypeMemory, duration)
	if err != nil {
		return nil, fmt.Errorf("获取内存指标失败: %v", err)
	}

	diskRecords, err := m.GetMetrics(MetricTypeDisk, duration)
	if err != nil {
		return nil, fmt.Errorf("获取磁盘指标失败: %v", err)
	}

	networkRecords, err := m.GetMetrics(MetricTypeNetwork, duration)
	if err != nil {
		return nil, fmt.Errorf("获取网络指标失败: %v", err)
	}

	// 构建系统概览（这里简化处理，实际应该聚合最新数据）
	overview := &SystemMetrics{
		CPU:     m.buildCPUMetrics(cpuRecords),
		Memory:  m.buildMemoryMetrics(memoryRecords),
		Disk:    m.buildDiskMetrics(diskRecords),
		Network: m.buildNetworkMetrics(networkRecords),
	}

	return overview, nil
}

// 构建CPU指标（简化实现）
func (m *Monitor) buildCPUMetrics(records []MetricRecord) CPUMetrics {
	if len(records) == 0 {
		return CPUMetrics{Timestamp: time.Now()}
	}

	// 取最新记录
	latest := records[len(records)-1]
	return CPUMetrics{
		Usage:     latest.Value,
		Timestamp: latest.Timestamp,
	}
}

// 构建内存指标（简化实现）
func (m *Monitor) buildMemoryMetrics(records []MetricRecord) MemoryMetrics {
	if len(records) == 0 {
		return MemoryMetrics{Timestamp: time.Now()}
	}

	latest := records[len(records)-1]
	return MemoryMetrics{
		Usage:     latest.Value,
		Timestamp: latest.Timestamp,
	}
}

// 构建磁盘指标（简化实现）
func (m *Monitor) buildDiskMetrics(records []MetricRecord) DiskMetrics {
	return DiskMetrics{
		Usage:     []DiskUsage{},
		Timestamp: time.Now(),
	}
}

// 构建网络指标（简化实现）
func (m *Monitor) buildNetworkMetrics(records []MetricRecord) NetworkMetrics {
	return NetworkMetrics{
		Interfaces: []NetworkInterface{},
		Timestamp:  time.Now(),
	}
}

// Cleanup 清理过期数据
func (m *Monitor) Cleanup() error {
	if m.storage == nil {
		return fmt.Errorf("存储未初始化")
	}

	return m.storage.Cleanup(m.ctx)
}
