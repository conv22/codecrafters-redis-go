package storage

type StorageValue struct {
	Value          string
	ExpirationTime *int64
}

type StorageKey struct {
	Key string
}

type StoragePair struct {
	Key   StorageKey
	Value StorageValue
}

type Storage interface {
	Get(key StorageKey) (StorageValue, bool)
	Set(key StorageKey, value StorageValue) error
	Delete(key StorageKey) error
}

func CreateStorage() Storage {
	return NewInMemoryStorage()
}
