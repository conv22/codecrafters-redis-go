package storage

type inMemoryStorage struct {
	data map[string]StorageValue
}

func NewInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		data: make(map[string]StorageValue),
	}
}

func (ims *inMemoryStorage) Get(key StorageKey) (StorageValue, bool) {
	value, ok := ims.data[key.Key]
	return value, ok
}

func (ims *inMemoryStorage) Set(key StorageKey, value StorageValue) error {
	ims.data[key.Key] = value
	return nil
}

func (ims *inMemoryStorage) Delete(key StorageKey) error {
	delete(ims.data, key.Key)
	return nil
}
