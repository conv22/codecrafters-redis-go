package storage

import (
	"testing"
)

func TestSetStorageById(t *testing.T) {
	collection := NewStorageCollection()
	storage := NewStorage()

	collection.SetStorageById(1, storage)

	if len(collection.Storages) != 1 {
		t.Errorf("Expected length of Storages map to be 1, got %d", len(collection.Storages))
	}
}

func TestGetCurrentStorage(t *testing.T) {
	collection := NewStorageCollection()
	storage := collection.GetCurrentStorage()

	if storage == nil {
		t.Error("Current storage should not be nil")
	}
}

func TestSetCurrentStorage(t *testing.T) {
	collection := NewStorageCollection()
	ims := collection.GetCurrentStorage()

	key := "testKey"

	err := ims.Set(key, &mockStorageValue)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	get, ok := ims.Get(key)

	if !ok {
		t.Errorf("No value returned")
	}

	if *get != mockStorageValue {
		t.Errorf("Expected value %v, got %v", mockStorageValue, get)
	}
}
