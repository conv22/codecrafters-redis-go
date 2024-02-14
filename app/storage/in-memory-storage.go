package storage

import (
	"path"
	"sync"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
)

type StorageKey = string

type StorageId = uint8

type StorageItem struct {
	Value    any
	Expiry   *time.Time
	Encoding string
}

type Storage struct {
	ID             uint8
	HashSize       int
	ExpireHashSize int
	CacheMap       map[StorageKey]StorageItem
	mu             sync.RWMutex
}

func NewStorage(config *config.Config) *Storage {
	if config != nil {
		savedDb := rdb.NewRdb()
		_, err := savedDb.HandleRead(path.Join(config.DirFlag, config.DbFilenameFlag))

		if err != nil {
			return &Storage{
				ID:       0,
				CacheMap: make(map[StorageKey]StorageItem),
			}
		}

		return &Storage{
			ID:       0,
			CacheMap: make(map[StorageKey]StorageItem),
		}

	} else {
		return &Storage{
			ID:       0,
			CacheMap: make(map[StorageKey]StorageItem),
		}
	}

}

func (ims *Storage) Get(key StorageKey) (StorageItem, bool) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, ok := ims.CacheMap[key]
	return value, ok
}

func (ims *Storage) GetKeys() []string {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	keys := make([]string, 0, len(ims.CacheMap))
	for k := range ims.CacheMap {
		keys = append(keys, k)
	}
	return keys
}

func (ims *Storage) Set(key StorageKey, value StorageItem) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.CacheMap[key] = value
	return nil
}

func (ims *Storage) Delete(key StorageKey) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.CacheMap, key)
	return nil
}
