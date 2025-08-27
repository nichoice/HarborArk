package monitor

import (
	"context"
	"time"
)

// MetricType 指标类型
type MetricType string

const (
	MetricTypeCPU      MetricType = "cpu"
	MetricTypeMemory   MetricType = "memory"
	MetricTypeDisk     MetricType = "disk"
	MetricTypeNetwork  MetricType = "network"
	MetricTypeNAS      MetricType = "nas"
	MetricTypeProtocol MetricType = "protocol"
)

// MetricRecord 指标记录
type MetricRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Timestamp  time.Time `gorm:"index" json:"timestamp"`
	MetricType string    `gorm:"index" json:"metric_type"`
	MetricName string    `gorm:"index" json:"metric_name"`
	Value      float64   `json:"value"`
	Labels     string    `gorm:"type:text" json:"labels"` // JSON格式存储标签
	CreatedAt  time.Time `json:"created_at"`
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPU     CPUMetrics     `json:"cpu"`
	Memory  MemoryMetrics  `json:"memory"`
	Disk    DiskMetrics    `json:"disk"`
	Network NetworkMetrics `json:"network"`
}

// CPUMetrics CPU指标
type CPUMetrics struct {
	Usage     float64   `json:"usage"`       // CPU使用率 (%)
	LoadAvg1  float64   `json:"load_avg_1"`  // 1分钟负载
	LoadAvg5  float64   `json:"load_avg_5"`  // 5分钟负载
	LoadAvg15 float64   `json:"load_avg_15"` // 15分钟负载
	Cores     int       `json:"cores"`       // CPU核心数
	Timestamp time.Time `json:"timestamp"`
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Total     uint64    `json:"total"`      // 总内存 (bytes)
	Used      uint64    `json:"used"`       // 已使用内存 (bytes)
	Available uint64    `json:"available"`  // 可用内存 (bytes)
	Usage     float64   `json:"usage"`      // 内存使用率 (%)
	SwapTotal uint64    `json:"swap_total"` // 交换分区总大小 (bytes)
	SwapUsed  uint64    `json:"swap_used"`  // 交换分区已使用 (bytes)
	Timestamp time.Time `json:"timestamp"`
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	Usage     []DiskUsage `json:"usage"` // 磁盘使用情况
	IO        DiskIO      `json:"io"`    // 磁盘I/O统计
	Timestamp time.Time   `json:"timestamp"`
}

// DiskUsage 磁盘使用情况
type DiskUsage struct {
	Device     string  `json:"device"`     // 设备名
	Mountpoint string  `json:"mountpoint"` // 挂载点
	Total      uint64  `json:"total"`      // 总大小 (bytes)
	Used       uint64  `json:"used"`       // 已使用 (bytes)
	Available  uint64  `json:"available"`  // 可用 (bytes)
	Usage      float64 `json:"usage"`      // 使用率 (%)
	Filesystem string  `json:"filesystem"` // 文件系统类型
}

// DiskIO 磁盘I/O统计
type DiskIO struct {
	ReadBytes  uint64 `json:"read_bytes"`  // 读取字节数
	WriteBytes uint64 `json:"write_bytes"` // 写入字节数
	ReadOps    uint64 `json:"read_ops"`    // 读取操作数
	WriteOps   uint64 `json:"write_ops"`   // 写入操作数
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	Timestamp  time.Time          `json:"timestamp"`
}

// NetworkInterface 网络接口
type NetworkInterface struct {
	Name        string `json:"name"`         // 接口名
	BytesRecv   uint64 `json:"bytes_recv"`   // 接收字节数
	BytesSent   uint64 `json:"bytes_sent"`   // 发送字节数
	PacketsRecv uint64 `json:"packets_recv"` // 接收包数
	PacketsSent uint64 `json:"packets_sent"` // 发送包数
	Errors      uint64 `json:"errors"`       // 错误数
	Drops       uint64 `json:"drops"`        // 丢包数
}

// NASMetrics NAS业务指标
type NASMetrics struct {
	ActiveConnections map[string]int `json:"active_connections"` // 活跃连接数
	TransferRate      TransferRate   `json:"transfer_rate"`      // 传输速率
	StorageUsage      []StorageUsage `json:"storage_usage"`      // 存储使用情况
	Timestamp         time.Time      `json:"timestamp"`
}

// TransferRate 传输速率
type TransferRate struct {
	Read  float64 `json:"read"`  // 读取速率 (MB/s)
	Write float64 `json:"write"` // 写入速率 (MB/s)
}

// StorageUsage 存储使用情况
type StorageUsage struct {
	Pool      string  `json:"pool"`      // 存储池名称
	Total     uint64  `json:"total"`     // 总容量 (bytes)
	Used      uint64  `json:"used"`      // 已使用 (bytes)
	Available uint64  `json:"available"` // 可用容量 (bytes)
	Usage     float64 `json:"usage"`     // 使用率 (%)
}

// ProtocolMetrics 协议服务指标
type ProtocolMetrics struct {
	ISCSI     ISCSIMetrics `json:"iscsi"`
	Samba     SambaMetrics `json:"samba"`
	NFS       NFSMetrics   `json:"nfs"`
	Timestamp time.Time    `json:"timestamp"`
}

// ISCSIMetrics iSCSI指标
type ISCSIMetrics struct {
	Targets   []ISCSITarget `json:"targets"`
	Sessions  int           `json:"sessions"`
	Connected bool          `json:"connected"`
}

// ISCSITarget iSCSI目标
type ISCSITarget struct {
	Name      string `json:"name"`
	LUN       int    `json:"lun"`
	Size      uint64 `json:"size"`
	Connected bool   `json:"connected"`
}

// SambaMetrics Samba指标
type SambaMetrics struct {
	Shares      []SambaShare `json:"shares"`
	Connections int          `json:"connections"`
	Running     bool         `json:"running"`
}

// SambaShare Samba共享
type SambaShare struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Connections int    `json:"connections"`
}

// NFSMetrics NFS指标
type NFSMetrics struct {
	Exports     []NFSExport `json:"exports"`
	Connections int         `json:"connections"`
	Running     bool        `json:"running"`
}

// NFSExport NFS导出
type NFSExport struct {
	Path        string `json:"path"`
	Clients     int    `json:"clients"`
	Permissions string `json:"permissions"`
}

// Collector 数据采集器接口
type Collector interface {
	// Collect 采集数据
	Collect(ctx context.Context) ([]MetricRecord, error)
	// GetType 获取采集器类型
	GetType() MetricType
	// IsEnabled 是否启用
	IsEnabled() bool
}

// Storage 存储接口
type Storage interface {
	// Store 存储指标数据
	Store(ctx context.Context, records []MetricRecord) error
	// Query 查询指标数据
	Query(ctx context.Context, req QueryRequest) ([]MetricRecord, error)
	// Cleanup 清理过期数据
	Cleanup(ctx context.Context) error
}

// QueryRequest 查询请求
type QueryRequest struct {
	MetricType string            `json:"metric_type"`
	MetricName string            `json:"metric_name"`
	StartTime  time.Time         `json:"start_time"`
	EndTime    time.Time         `json:"end_time"`
	Labels     map[string]string `json:"labels"`
	Limit      int               `json:"limit"`
}

// Scheduler 调度器接口
type Scheduler interface {
	// Start 启动调度器
	Start(ctx context.Context, collectFunc func() error) error
	// Stop 停止调度器
	Stop() error
	// SetInterval 设置采集间隔
	SetInterval(interval time.Duration)
}
