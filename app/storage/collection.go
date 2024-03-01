package storage

import "sync"

type StorageCollection struct {
	CurrStorageId storageId
	Storages      map[storageId]*Storage
	AuxFields     map[string]interface{}
	mu            sync.RWMutex
}

func NewStorageCollection() *StorageCollection {
	return &StorageCollection{
		Storages:  make(map[storageId]*Storage),
		AuxFields: make(map[string]interface{}),
	}

}

func (collection *StorageCollection) SetAuxField(key string, value interface{}) error {
	collection.mu.Lock()
	defer collection.mu.Unlock()
	collection.AuxFields[key] = value
	return nil
}

func (collection *StorageCollection) SetStorageById(id storageId, storage *Storage) {
	storage.mu.Lock()
	defer storage.mu.Unlock()
	collection.Storages[id] = storage
}

func (collection *StorageCollection) GetCurrentStorage() *Storage {
	collection.mu.Lock()
	defer collection.mu.Unlock()
	storage, ok := collection.Storages[collection.CurrStorageId]

	if !ok {
		newStorage := NewStorage()
		collection.SetStorageById(collection.CurrStorageId, newStorage)
		return newStorage
	}

	return storage
}

func (collection *StorageCollection) SetItemToCurrentStorage(key storageKey, item *StorageItem) error {
	storage := collection.GetCurrentStorage()

	err := storage.Set(key, item)

	if err != nil {
		return err
	}
	return nil
}

func (collection *StorageCollection) GetItemFromCurrentStorage(key storageKey) (*StorageItem, bool) {
	storage := collection.GetCurrentStorage()

	item, ok := storage.Get(key)

	return item, ok
}
