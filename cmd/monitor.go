package cmd

import (
	"HarborArk/config"
	"HarborArk/internal/monitor"
	"HarborArk/internal/monitor/collector"
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "显示系统监控信息",
	Long:  "在终端中显示实时系统监控信息",
	Run: func(cmd *cobra.Command, args []string) {
		runMonitor(cmd, args)
	},
}

func init() {
	monitorCmd.Flags().StringP("type", "t", "system", "监控类型 (system, cpu, memory, disk, network)")
	monitorCmd.Flags().StringP("interval", "i", "2s", "刷新间隔")
	monitorCmd.Flags().IntP("count", "c", 0, "显示次数 (0表示持续显示)")
	monitorCmd.Flags().BoolP("once", "o", false, "只显示一次")

	rootCmd.AddCommand(monitorCmd)
}

func runMonitor(cmd *cobra.Command, args []string) {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		return
	}

	// 初始化数据库
	config.ConnectDB()

	// 获取参数
	monitorType, _ := cmd.Flags().GetString("type")
	intervalStr, _ := cmd.Flags().GetString("interval")
	count, _ := cmd.Flags().GetInt("count")
	once, _ := cmd.Flags().GetBool("once")

	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		fmt.Printf("无效的间隔时间: %v\n", err)
		return
	}

	// 创建采集器
	systemCollector := collector.NewSystemCollector()

	fmt.Printf("HarborArk 系统监控 - 类型: %s, 间隔: %s\n", monitorType, intervalStr)
	fmt.Println("按 Ctrl+C 退出")
	fmt.Println("--------------------------------------------------------------------------------")

	iteration := 0
	for {
		// 采集数据
		ctx := context.Background()
		records, err := systemCollector.Collect(ctx)
		if err != nil {
			fmt.Printf("采集数据失败: %v\n", err)
			time.Sleep(interval)
			continue
		}

		// 显示数据
		displayMetrics(records, monitorType)

		iteration++

		// 检查是否只显示一次
		if once || (count > 0 && iteration >= count) {
			break
		}

		time.Sleep(interval)
	}
}

func displayMetrics(records []monitor.MetricRecord, monitorType string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("\n[%s] 系统监控数据:\n", timestamp)

	// 按类型分组显示
	cpuMetrics := make(map[string]float64)
	memoryMetrics := make(map[string]float64)
	diskMetrics := make(map[string]float64)
	networkMetrics := make(map[string]float64)

	for _, record := range records {
		switch record.MetricType {
		case "cpu":
			cpuMetrics[record.MetricName] = record.Value
		case "memory":
			memoryMetrics[record.MetricName] = record.Value
		case "disk":
			diskMetrics[record.MetricName] = record.Value
		case "network":
			networkMetrics[record.MetricName] = record.Value
		}
	}

	// 显示CPU信息
	if monitorType == "system" || monitorType == "cpu" {
		fmt.Println("CPU:")
		if usage, ok := cpuMetrics["cpu_usage_percent"]; ok {
			fmt.Printf("  使用率: %.1f%%\n", usage)
		}
		if cores, ok := cpuMetrics["cpu_cores"]; ok {
			fmt.Printf("  核心数: %.0f\n", cores)
		}
		if load1, ok := cpuMetrics["load_avg_1"]; ok {
			fmt.Printf("  负载: %.2f", load1)
			if load5, ok := cpuMetrics["load_avg_5"]; ok {
				fmt.Printf(", %.2f", load5)
				if load15, ok := cpuMetrics["load_avg_15"]; ok {
					fmt.Printf(", %.2f", load15)
				}
			}
			fmt.Println()
		}
	}

	// 显示内存信息
	if monitorType == "system" || monitorType == "memory" {
		fmt.Println("内存:")
		if usage, ok := memoryMetrics["memory_usage_percent"]; ok {
			fmt.Printf("  使用率: %.1f%%\n", usage)
		}
		if total, ok := memoryMetrics["memory_total_bytes"]; ok {
			fmt.Printf("  总量: %s\n", formatBytes(uint64(total)))
		}
		if used, ok := memoryMetrics["memory_used_bytes"]; ok {
			fmt.Printf("  已使用: %s\n", formatBytes(uint64(used)))
		}
		if available, ok := memoryMetrics["memory_available_bytes"]; ok {
			fmt.Printf("  可用: %s\n", formatBytes(uint64(available)))
		}
	}

	// 显示磁盘信息
	if monitorType == "system" || monitorType == "disk" {
		fmt.Println("磁盘:")
		diskCount := 0
		for name, value := range diskMetrics {
			if name == "disk_usage_percent" {
				diskCount++
				fmt.Printf("  使用率: %.1f%%\n", value)
			}
		}
		if diskCount == 0 {
			fmt.Println("  暂无磁盘数据")
		}
	}

	// 显示网络信息
	if monitorType == "system" || monitorType == "network" {
		fmt.Println("网络:")
		networkCount := 0
		for name, value := range networkMetrics {
			if name == "network_receive_bytes" {
				networkCount++
				fmt.Printf("  接收: %s\n", formatBytes(uint64(value)))
			} else if name == "network_transmit_bytes" {
				fmt.Printf("  发送: %s\n", formatBytes(uint64(value)))
			}
		}
		if networkCount == 0 {
			fmt.Println("  暂无网络数据")
		}
	}
}

// formatBytes 格式化字节数
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// MemoryStorage 内存存储实现（用于CLI测试）
type MemoryStorage struct {
	records []monitor.MetricRecord
}

func (s *MemoryStorage) Store(ctx context.Context, records []monitor.MetricRecord) error {
	s.records = append(s.records, records...)
	return nil
}

func (s *MemoryStorage) Query(ctx context.Context, req monitor.QueryRequest) ([]monitor.MetricRecord, error) {
	return s.records, nil
}

func (s *MemoryStorage) Cleanup(ctx context.Context) error {
	s.records = s.records[:0]
	return nil
}
