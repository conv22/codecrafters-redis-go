package storage

import "sync"

type StorageCollection struct {
	CurrStorageId StorageId
	Storages      map[StorageId]*Storage
	AuxFields     map[string]interface{}
	mu            sync.RWMutex
}

func NewStorageCollection() *StorageCollection {
	return &StorageCollection{
		Storages:  make(map[StorageId]*Storage),
		AuxFields: make(map[string]interface{}),
	}

}

func (collection *StorageCollection) SetAuxField(key string, value interface{}) error {
	collection.mu.Lock()
	defer collection.mu.Unlock()
	collection.AuxFields[key] = value
	return nil
}

func (collection *StorageCollection) SetStorageById(id StorageId, storage *Storage) {
	storage.mu.Lock()
	collection.Storages[id] = storage
	storage.mu.Unlock()
}

func (collection *StorageCollection) GetCurrentStorage() *Storage {
	storage, ok := collection.Storages[collection.CurrStorageId]

	if !ok {
		newStorage := NewStorage()
		collection.SetStorageById(collection.CurrStorageId, newStorage)
		return newStorage
	}

	return storage
}

func (collection *StorageCollection) SetItemToCurrentStorage(key StorageKey, item *StorageItem) error {
	storage := collection.GetCurrentStorage()

	err := storage.Set(key, item)

	if err != nil {
		return err
	}
	return nil
}

func (collection *StorageCollection) GetItemFromCurrentStorage(key StorageKey) (*StorageItem, bool) {
	storage := collection.GetCurrentStorage()

	item, ok := storage.Get(key)

	return item, ok
}
