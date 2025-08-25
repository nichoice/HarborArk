package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
)

var badgerDB *badger.DB

// InitBadger 初始化 BadgerDB，并启动定期 GC
func InitBadger() error {
	cfg := GetBadgerConfig()
	if cfg.Path == "" {
		cfg.Path = "data/badger"
	}
	if err := os.MkdirAll(cfg.Path, 0o755); err != nil {
		return err
	}

	opts := badger.DefaultOptions(filepath.Clean(cfg.Path))
	// 降低 Badger 自身日志噪音
	opts = opts.WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	badgerDB = db
	zap.L().Info("BadgerDB 初始化成功", zap.String("path", cfg.Path))

	// 启动后台 GC
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if badgerDB == nil {
				return
			}
			// 重复多次尝试回收
			for i := 0; i < 3; i++ {
				if err := badgerDB.RunValueLogGC(0.5); err != nil {
					break
				}
			}
		}
	}()

	return nil
}

// GetBadger 获取 BadgerDB 实例
func GetBadger() *badger.DB {
	return badgerDB
}

// CloseBadger 关闭 BadgerDB
func CloseBadger() error {
	if badgerDB != nil {
		return badgerDB.Close()
	}
	return nil
}