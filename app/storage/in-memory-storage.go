package storage

import "sync"

type inMemoryStorage struct {
	data map[string]StorageValue
	mu   sync.RWMutex
}

func NewInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		data: make(map[string]StorageValue),
	}
}

func (ims *inMemoryStorage) Get(key StorageKey) (StorageValue, bool) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, ok := ims.data[key.Key]
	return value, ok
}

func (ims *inMemoryStorage) Set(key StorageKey, value StorageValue) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.data[key.Key] = value
	return nil
}

func (ims *inMemoryStorage) Delete(key StorageKey) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.data, key.Key)
	return nil
}
