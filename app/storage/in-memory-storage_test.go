package storage

import (
	"testing"
)

var mockStorageValue = StorageValue{Value: "testValue"}

func TestGet(t *testing.T) {
	ims := NewInMemoryStorage()

	key := StorageKey{Key: "testKey"}
	ims.Set(key, mockStorageValue)

	result := ims.Get(key)

	if *result != mockStorageValue {
		t.Errorf("Expected value %v, got %v", mockStorageValue, result)
	}

	result = ims.Get(StorageKey{Key: "nonExistingKey"})
	if result != nil {
		t.Errorf("Expected error 'key not found', got %v", *result)

	}

}

func TestSet(t *testing.T) {
	ims := NewInMemoryStorage()

	key := StorageKey{Key: "testKey"}

	err := ims.Set(key, mockStorageValue)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	get := ims.Get(key)

	if get == nil {
		t.Errorf("No value returned")
	}

	if *get != mockStorageValue {
		t.Errorf("Expected value %v, got %v", mockStorageValue, get)
	}
}
