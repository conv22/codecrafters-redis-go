package storage

import "sync"

type StorageValue struct {
	Value          any
	ExpirationTime *int64
}

type StorageKey struct {
	Key string
}

type InMemoryStorage struct {
	data map[string]StorageValue
	mu   sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]StorageValue),
	}
}

func (ims *InMemoryStorage) Get(key StorageKey) (StorageValue, bool) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, ok := ims.data[key.Key]
	return value, ok
}

func (ims *InMemoryStorage) Set(key StorageKey, value StorageValue) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.data[key.Key] = value
	return nil
}

func (ims *InMemoryStorage) Delete(key StorageKey) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.data, key.Key)
	return nil
}
