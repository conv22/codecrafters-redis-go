package storage

type StorageCollection struct {
	CurrStorageId StorageId
	Storages      map[StorageId]*Storage
}

func NewStorageCollection() *StorageCollection {
	return &StorageCollection{
		Storages: make(map[StorageId]*Storage),
	}

}

func (collection *StorageCollection) SetStorageById(id StorageId, storage *Storage) {
	collection.Storages[id] = storage
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
