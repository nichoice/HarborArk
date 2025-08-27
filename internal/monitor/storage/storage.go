package storage

import (
	"HarborArk/internal/monitor"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DBStorage 数据库存储实现
type DBStorage struct {
	db *gorm.DB
}

// NewDBStorage 创建数据库存储
func NewDBStorage(db *gorm.DB) *DBStorage {
	return &DBStorage{
		db: db,
	}
}

// Store 存储指标数据
func (s *DBStorage) Store(ctx context.Context, records []monitor.MetricRecord) error {
	if len(records) == 0 {
		return nil
	}

	// 批量插入
	if err := s.db.WithContext(ctx).CreateInBatches(records, 100).Error; err != nil {
		return fmt.Errorf("批量插入指标数据失败: %v", err)
	}

	zap.L().Debug("存储指标数据成功", zap.Int("count", len(records)))
	return nil
}

// Query 查询指标数据
func (s *DBStorage) Query(ctx context.Context, req monitor.QueryRequest) ([]monitor.MetricRecord, error) {
	var records []monitor.MetricRecord

	query := s.db.WithContext(ctx).Model(&monitor.MetricRecord{})

	// 添加查询条件
	if req.MetricType != "" {
		query = query.Where("metric_type = ?", req.MetricType)
	}

	if req.MetricName != "" {
		query = query.Where("metric_name = ?", req.MetricName)
	}

	if !req.StartTime.IsZero() {
		query = query.Where("timestamp >= ?", req.StartTime)
	}

	if !req.EndTime.IsZero() {
		query = query.Where("timestamp <= ?", req.EndTime)
	}

	// 标签过滤
	if len(req.Labels) > 0 {
		for key, value := range req.Labels {
			query = query.Where("JSON_EXTRACT(labels, ?) = ?", "$."+key, value)
		}
	}

	// 排序和限制
	query = query.Order("timestamp DESC")
	if req.Limit > 0 {
		query = query.Limit(req.Limit)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("查询指标数据失败: %v", err)
	}

	return records, nil
}

// Cleanup 清理过期数据
func (s *DBStorage) Cleanup(ctx context.Context) error {
	// 删除7天前的原始数据
	cutoffTime := time.Now().AddDate(0, 0, -7)

	result := s.db.WithContext(ctx).
		Where("created_at < ?", cutoffTime).
		Delete(&monitor.MetricRecord{})

	if result.Error != nil {
		return fmt.Errorf("清理过期数据失败: %v", result.Error)
	}

	if result.RowsAffected > 0 {
		zap.L().Info("清理过期数据完成", zap.Int64("deleted", result.RowsAffected))
	}

	return nil
}

// GetLatestMetrics 获取最新指标
func (s *DBStorage) GetLatestMetrics(ctx context.Context, metricType string, limit int) ([]monitor.MetricRecord, error) {
	var records []monitor.MetricRecord

	query := s.db.WithContext(ctx).
		Where("metric_type = ?", metricType).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("获取最新指标失败: %v", err)
	}

	return records, nil
}

// GetMetricsByTimeRange 按时间范围获取指标
func (s *DBStorage) GetMetricsByTimeRange(ctx context.Context, metricType, metricName string, start, end time.Time) ([]monitor.MetricRecord, error) {
	var records []monitor.MetricRecord

	query := s.db.WithContext(ctx).
		Where("metric_type = ? AND timestamp BETWEEN ? AND ?", metricType, start, end)

	if metricName != "" {
		query = query.Where("metric_name = ?", metricName)
	}

	if err := query.Order("timestamp ASC").Find(&records).Error; err != nil {
		return nil, fmt.Errorf("按时间范围查询指标失败: %v", err)
	}

	return records, nil
}

// AggregateMetrics 聚合指标数据
func (s *DBStorage) AggregateMetrics(ctx context.Context, metricType string, interval time.Duration, start, end time.Time) ([]AggregatedMetric, error) {
	var results []AggregatedMetric

	// 根据间隔确定时间格式
	var timeFormat string
	switch {
	case interval >= 24*time.Hour:
		timeFormat = "%Y-%m-%d"
	case interval >= time.Hour:
		timeFormat = "%Y-%m-%d %H:00:00"
	default:
		timeFormat = "%Y-%m-%d %H:%i:00"
	}

	query := `
		SELECT 
			metric_name,
			DATE_FORMAT(timestamp, ?) as time_bucket,
			AVG(value) as avg_value,
			MIN(value) as min_value,
			MAX(value) as max_value,
			COUNT(*) as count
		FROM metric_records 
		WHERE metric_type = ? AND timestamp BETWEEN ? AND ?
		GROUP BY metric_name, time_bucket
		ORDER BY time_bucket ASC
	`

	if err := s.db.WithContext(ctx).Raw(query, timeFormat, metricType, start, end).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("聚合指标数据失败: %v", err)
	}

	return results, nil
}

// AggregatedMetric 聚合指标结果
type AggregatedMetric struct {
	MetricName string    `json:"metric_name"`
	TimeBucket string    `json:"time_bucket"`
	AvgValue   float64   `json:"avg_value"`
	MinValue   float64   `json:"min_value"`
	MaxValue   float64   `json:"max_value"`
	Count      int       `json:"count"`
	Timestamp  time.Time `json:"timestamp"`
}

// CreateMetricRecord 创建指标记录的辅助函数
func CreateMetricRecord(metricType, metricName string, value float64, labels map[string]string) monitor.MetricRecord {
	var labelsJSON string
	if len(labels) > 0 {
		if data, err := json.Marshal(labels); err == nil {
			labelsJSON = string(data)
		}
	}

	return monitor.MetricRecord{
		Timestamp:  time.Now(),
		MetricType: metricType,
		MetricName: metricName,
		Value:      value,
		Labels:     labelsJSON,
		CreatedAt:  time.Now(),
	}
}
