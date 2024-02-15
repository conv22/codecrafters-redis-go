package storage

import (
	"sync"
	"time"
)

type StorageKey = string

type StorageId = uint8

type StorageItem struct {
	Value    any
	Expiry   *time.Time
	Encoding byte
}

type Storage struct {
	HashSize       int
	ExpireHashSize int
	CacheMap       map[StorageKey]*StorageItem
	mu             sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		CacheMap: make(map[StorageKey]*StorageItem),
	}
}

func (ims *Storage) Get(key StorageKey) (*StorageItem, bool) {
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

func (ims *Storage) Set(key StorageKey, value *StorageItem) error {
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
