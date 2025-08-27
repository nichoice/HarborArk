package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

// DefaultScheduler 默认调度器实现
type DefaultScheduler struct {
	interval time.Duration
	ticker   *time.Ticker
	running  bool
	mu       sync.RWMutex
	stopCh   chan struct{}
}

// NewDefaultScheduler 创建默认调度器
func NewDefaultScheduler() *DefaultScheduler {
	return &DefaultScheduler{
		interval: 15 * time.Second, // 默认15秒
		stopCh:   make(chan struct{}),
	}
}

// SetInterval 设置采集间隔
func (s *DefaultScheduler) SetInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.interval = interval

	// 如果正在运行，重新创建ticker
	if s.running && s.ticker != nil {
		s.ticker.Stop()
		s.ticker = time.NewTicker(interval)
	}
}

// Start 启动调度器
func (s *DefaultScheduler) Start(ctx context.Context, collectFunc func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("调度器已在运行")
	}

	s.ticker = time.NewTicker(s.interval)
	s.running = true

	go s.run(ctx, collectFunc)

	zap.L().Info("调度器启动成功", zap.Duration("interval", s.interval))
	return nil
}

// Stop 停止调度器
func (s *DefaultScheduler) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	close(s.stopCh)

	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}

	s.running = false

	zap.L().Info("调度器已停止")
	return nil
}

// run 运行调度循环
func (s *DefaultScheduler) run(ctx context.Context, collectFunc func() error) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Error("调度器运行时发生panic", zap.Any("panic", r))
		}
	}()

	// 立即执行一次采集
	if err := collectFunc(); err != nil {
		zap.L().Error("初始数据采集失败", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			zap.L().Info("调度器收到上下文取消信号")
			return

		case <-s.stopCh:
			zap.L().Info("调度器收到停止信号")
			return

		case <-s.ticker.C:
			start := time.Now()

			if err := collectFunc(); err != nil {
				zap.L().Error("定时数据采集失败", zap.Error(err))
			}

			duration := time.Since(start)
			if duration > s.interval/2 {
				zap.L().Warn("数据采集耗时过长",
					zap.Duration("duration", duration),
					zap.Duration("interval", s.interval))
			}
		}
	}
}

// IsRunning 检查是否运行中
func (s *DefaultScheduler) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
