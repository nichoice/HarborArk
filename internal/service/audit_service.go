package service

import (
	"HarborArk/config"
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
)

const auditPrefix = "audit:"

type AuditRecord struct {
	Timestamp  int64                  `json:"timestamp"`   // Unix ms
	ActorID    uint                   `json:"actor_id"`    // 操作者ID
	ActorName  string                 `json:"actor_name"`  // 操作者名称
	Action     string                 `json:"action"`      // 操作类型：delete/rename/mkdir/update_meta/...
	TargetPath string                 `json:"target_path"` // 目标路径
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// LogAudit 记录审计事件
func LogAudit(ctx context.Context, rec AuditRecord) error {
	db := config.GetBadger()
	if db == nil {
		return fmt.Errorf("badger not initialized")
	}
	if rec.Timestamp == 0 {
		rec.Timestamp = time.Now().UnixMilli()
	}

	key := makeAuditKey(rec.Timestamp)
	val, _ := json.Marshal(rec)

	return db.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), val)
		return txn.SetEntry(e)
	})
}

func makeAuditKey(ts int64) string {
	suffix := make([]byte, 8)
	_, _ = rand.Read(suffix)
	return fmt.Sprintf("%s%d:%s", auditPrefix, ts, hex.EncodeToString(suffix))
}

func parseAuditKey(key []byte) (int64, bool) {
	s := string(key)
	if !strings.HasPrefix(s, auditPrefix) {
		return 0, false
	}
	s = strings.TrimPrefix(s, auditPrefix)
	parts := strings.SplitN(s, ":", 2)
	if len(parts) < 1 {
		return 0, false
	}
	t, err := timeFromString(parts[0])
	if err != nil {
		return 0, false
	}
	return t, true
}

func timeFromString(s string) (int64, error) {
	var t int64
	_, err := fmt.Sscanf(s, "%d", &t)
	return t, err
}

// StartAuditRetentionWorker 启动审计日志保留清理任务（按天保留）
func StartAuditRetentionWorker(retentionDays int) {
	if retentionDays <= 0 {
		retentionDays = 30
	}
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour).UnixMilli()
			err := purgeAuditBefore(cutoff)
			if err != nil {
				zap.L().Warn("审计日志清理失败", zap.Error(err))
			}
		}
	}()
}

func purgeAuditBefore(cutoffMs int64) error {
	db := config.GetBadger()
	if db == nil {
		return fmt.Errorf("badger not initialized")
	}

	// 使用 WriteBatch 批量删除，避免单个事务过大
	wb := db.NewWriteBatch()
	defer wb.Cancel()

	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Prefix:         []byte(auditPrefix),
		})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			ts, ok := parseAuditKey(item.Key())
			if !ok || ts >= cutoffMs {
				continue
			}
			k := append([]byte(nil), item.Key()...)
			if err := wb.Delete(k); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return wb.Flush()
}

// ExportAudit 导出[from, to]时间范围内的审计日志，生成zip归档，返回归档文件路径
func ExportAudit(from, to time.Time, exportDir string) (string, error) {
	db := config.GetBadger()
	if db == nil {
		return "", fmt.Errorf("badger not initialized")
	}
	if exportDir == "" {
		exportDir = "exports/audit"
	}
	if err := os.MkdirAll(exportDir, 0o755); err != nil {
		return "", err
	}

	jsonlName := fmt.Sprintf("audit_%s_%s.jsonl",
		from.Format("20060102_150405"),
		to.Format("20060102_150405"),
	)
	zipName := strings.TrimSuffix(jsonlName, ".jsonl") + ".zip"
	jsonlPath := filepath.Join(exportDir, jsonlName)
	zipPath := filepath.Join(exportDir, zipName)

	// 写入 JSONL 到内存 buffer，最后写入 zip
	var buf bytes.Buffer

	fromMs := from.UnixMilli()
	toMs := to.UnixMilli()

	err := db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: true,
			Prefix:         []byte(auditPrefix),
		})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			ts, ok := parseAuditKey(item.Key())
			if !ok || ts < fromMs || ts > toMs {
				continue
			}
			err := item.Value(func(v []byte) error {
				buf.Write(v)
				buf.WriteByte('\n')
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	// 同时落地一份 JSONL
	if err := os.WriteFile(jsonlPath, buf.Bytes(), 0o644); err != nil {
		return "", err
	}

	// 写 zip 归档
	zf, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}
	defer zf.Close()

	zw := zip.NewWriter(zf)
	defer zw.Close()

	w, err := zw.Create(jsonlName)
	if err != nil {
		return "", err
	}
	if _, err := w.Write(buf.Bytes()); err != nil {
		return "", err
	}

	if err := zw.Close(); err != nil {
		return "", err
	}

	zap.L().Info("审计日志导出完成", zap.String("zip", zipPath), zap.String("jsonl", jsonlPath))
	return zipPath, nil
}
