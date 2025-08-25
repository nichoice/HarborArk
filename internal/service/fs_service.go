package service

import (
	"HarborArk/config"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
)

const (
	metaPrefix = "meta:"
)

// FileMetadata 仅存放元数据（基础信息不入库）
type FileMetadata struct {
	Tags          []string          `json:"tags,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	Custom        map[string]string `json:"custom,omitempty"`
	CreatedBy     uint              `json:"created_by,omitempty"`
	CreatedByName string            `json:"created_by_name,omitempty"`
	CreatedAt     int64             `json:"created_at,omitempty"` // ms
	UpdatedBy     uint              `json:"updated_by,omitempty"`
	UpdatedByName string            `json:"updated_by_name,omitempty"`
	UpdatedAt     int64             `json:"updated_at,omitempty"` // ms
}

// FileItem 返回列表条目（基础信息来自文件系统）
type FileItem struct {
	Name    string        `json:"name"`
	Path    string        `json:"path"`
	IsDir   bool          `json:"is_dir"`
	Size    int64         `json:"size"`
	Mode    uint32        `json:"mode"`
	ModTime int64         `json:"mod_time"` // ms
	Meta    *FileMetadata `json:"meta,omitempty"`
}

type FSService struct{}

// validatePath 验证路径是否在允许的目录范围内
func validatePath(path string) error {
	fmConfig := config.GetFileManagerConfig()
	if !fmConfig.RestrictToAllowedDirs {
		return nil
	}

	cleanPath := filepath.Clean(path)
	for _, allowedDir := range fmConfig.AllowedDirs {
		allowedClean := filepath.Clean(allowedDir)
		if strings.HasPrefix(cleanPath, allowedClean) {
			return nil
		}
	}

	return errors.New("访问被拒绝：路径不在允许的目录范围内")
}

// ListDir 以分页方式列出目录，避免一次性加载巨量文件
// 返回 items 与 hasMore（是否还有后续数据）
func (s *FSService) ListDir(p string, offset, limit int, includeHidden bool) ([]FileItem, bool, error) {
	// 验证路径权限
	if err := validatePath(p); err != nil {
		return nil, false, err
	}
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	cleanPath := filepath.Clean(p)
	f, err := os.Open(cleanPath)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()

	// 预读 offset+limit，但采用分批次读取避免占用过大内存
	names := make([]string, 0, limit)
	skipped := 0
	const chunk = 1024
	hasMore := false

	for len(names) < limit {
		toRead := chunk
		if offset+limit > 0 && toRead < limit {
			// 使用固定 chunk 即可
		}
		batch, err := f.Readdirnames(toRead)
		if err != nil && err != io.EOF {
			return nil, false, err
		}
		if len(batch) == 0 {
			break
		}
		for _, n := range batch {
			// 过滤隐藏文件
			if !includeHidden && strings.HasPrefix(n, ".") {
				if skipped < offset {
					skipped++
				}
				continue
			}
			if skipped < offset {
				skipped++
				continue
			}
			if len(names) < limit {
				names = append(names, n)
			} else {
				hasMore = true
				// 仍需继续消费后续 batch 才能确定 hasMore，既然到这已确定，直接标记后返回
				// 但为了减少系统调用，提前退出外层循环
				goto BUILD
			}
		}
		if err == io.EOF {
			break
		}
	}
BUILD:
	// 可选：对名称排序（按名称升序），避免不同文件系统返回顺序不稳定
	sort.Strings(names)

	// 组装 FileItem，并从 Badger 读取元数据
	items := make([]FileItem, 0, len(names))
	db := config.GetBadger()
	var viewTxn *badger.Txn
	if db != nil {
		viewTxn = db.NewTransaction(false)
		defer viewTxn.Discard()
	}

	for _, name := range names {
		full := filepath.Join(cleanPath, name)
		fi, statErr := os.Lstat(full)
		if statErr != nil {
			// 跳过无法 stat 的项
			continue
		}
		item := FileItem{
			Name:    name,
			Path:    full,
			IsDir:   fi.IsDir(),
			Size:    sizeOf(fi),
			Mode:    uint32(fi.Mode().Perm()),
			ModTime: fi.ModTime().UnixMilli(),
		}

		if viewTxn != nil {
			if meta, _ := getMetaTxn(viewTxn, full); meta != nil {
				item.Meta = meta
			}
		}
		items = append(items, item)
	}

	return items, hasMore, nil
}

func sizeOf(fi os.FileInfo) int64 {
	if fi.IsDir() {
		return 0
	}
	return fi.Size()
}

func metaKey(path string) []byte {
	return []byte(metaPrefix + filepath.Clean(path))
}

func getMetaTxn(txn *badger.Txn, path string) (*FileMetadata, error) {
	item, err := txn.Get(metaKey(path))
	if errors.Is(err, badger.ErrKeyNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var meta FileMetadata
	if err := item.Value(func(v []byte) error {
		return json.Unmarshal(v, &meta)
	}); err != nil {
		return nil, err
	}
	return &meta, nil
}

func GetMetadata(path string) (*FileMetadata, error) {
	// 验证路径权限
	if err := validatePath(path); err != nil {
		return nil, err
	}

	db := config.GetBadger()
	if db == nil {
		return nil, nil
	}
	var meta *FileMetadata
	err := db.View(func(txn *badger.Txn) error {
		m, err := getMetaTxn(txn, path)
		if err != nil {
			return err
		}
		meta = m
		return nil
	})
	return meta, err
}

func UpsertMetadata(path string, meta FileMetadata) error {
	// 验证路径权限
	if err := validatePath(path); err != nil {
		return err
	}

	db := config.GetBadger()
	if db == nil {
		return nil
	}
	now := time.Now().UnixMilli()
	if meta.CreatedAt == 0 {
		meta.CreatedAt = now
	}
	meta.UpdatedAt = now

	val, _ := json.Marshal(meta)
	return db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(badger.NewEntry(metaKey(path), val))
	})
}

func RenamePath(oldPath, newPath string) error {
	// 验证路径权限
	if err := validatePath(oldPath); err != nil {
		return err
	}
	if err := validatePath(newPath); err != nil {
		return err
	}

	// 先重命名文件系统
	if err := os.Rename(filepath.Clean(oldPath), filepath.Clean(newPath)); err != nil {
		return err
	}

	// 迁移元数据键
	db := config.GetBadger()
	if db == nil {
		return nil
	}
	oldK := metaKey(oldPath)
	newK := metaKey(newPath)

	return db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(oldK)
		if err == nil {
			var v []byte
			if err := item.Value(func(b []byte) error {
				v = append([]byte(nil), b...)
				return nil
			}); err != nil {
				return err
			}
			if err := txn.Set(newK, v); err != nil {
				return err
			}
			if err := txn.Delete(oldK); err != nil {
				return err
			}
		}
		// 若是目录，迁移子路径元数据（前缀替换）
		fi, _ := os.Lstat(newPath)
		if fi != nil && fi.IsDir() {
			prefixOld := metaPrefix + filepath.Clean(oldPath) + string(os.PathSeparator)
			prefixNew := metaPrefix + filepath.Clean(newPath) + string(os.PathSeparator)
			it := txn.NewIterator(badger.IteratorOptions{
				PrefetchValues: true,
				Prefix:         []byte(prefixOld),
			})
			defer it.Close()
			wb := txn
			for it.Rewind(); it.Valid(); it.Next() {
				key := append([]byte(nil), it.Item().Key()...)
				newKey := []byte(strings.Replace(string(key), prefixOld, prefixNew, 1))
				var v []byte
				if err := it.Item().Value(func(b []byte) error {
					v = append([]byte(nil), b...)
					return nil
				}); err != nil {
					return err
				}
				if err := wb.Set(newKey, v); err != nil {
					return err
				}
				if err := wb.Delete(key); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func DeletePath(path string) error {
	// 验证路径权限
	if err := validatePath(path); err != nil {
		return err
	}

	clean := filepath.Clean(path)
	// 删除文件/目录（目录递归）
	if err := os.RemoveAll(clean); err != nil {
		return err
	}

	// 清理元数据
	db := config.GetBadger()
	if db == nil {
		return nil
	}
	return db.Update(func(txn *badger.Txn) error {
		// 删除自身
		_ = txn.Delete(metaKey(clean))

		// 若目录，删除所有子项元数据
		prefix := metaPrefix + clean + string(os.PathSeparator)
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Prefix:         []byte(prefix),
		})
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			if err := txn.Delete(it.Item().Key()); err != nil {
				return err
			}
		}
		return nil
	})
}

func Mkdir(parent, name string) error {
	// 验证路径权限
	if err := validatePath(parent); err != nil {
		return err
	}

	p := filepath.Join(filepath.Clean(parent), name)

	// 验证新创建的目录路径
	if err := validatePath(p); err != nil {
		return err
	}

	return os.MkdirAll(p, 0o755)
}

// Audit wrappers

func LogFSAudit(ctx context.Context, actorID uint, actorName, action, target string, extra map[string]interface{}) {
	_ = LogAudit(ctx, AuditRecord{
		ActorID:    actorID,
		ActorName:  actorName,
		Action:     action,
		TargetPath: filepath.Clean(target),
		Extra:      extra,
	})
}
