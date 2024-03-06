package storage

import (
	"sync"
)

type storageKey = string

type storageId = uint8

const (
	STRING_TYPE = "string"
	NONE_TYPE   = "none"
	STREAM_TYPE = "stream"
)

type StorageItem struct {
	Value    any
	Type     string
	ExpiryMs int64
	Encoding byte
}

type Storage struct {
	HashSize       int
	ExpireHashSize int
	CacheMap       map[storageKey]*StorageItem
	mu             sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		CacheMap: make(map[storageKey]*StorageItem),
	}
}

func (ims *Storage) Get(key storageKey) (*StorageItem, bool) {
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

func (ims *Storage) Set(key storageKey, value *StorageItem) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.CacheMap[key] = value
	return nil
}

func (ims *Storage) Delete(key storageKey) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.CacheMap, key)
	return nil
}
