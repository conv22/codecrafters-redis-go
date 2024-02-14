package storage

import (
	"testing"
)

var mockStorageValue = StorageItem{Value: "testValue"}

func TestGet(t *testing.T) {
	ims := NewStorage()

	key := "testKey"
	ims.Set(key, &mockStorageValue)

	result, ok := ims.Get(key)

	if !ok || *result != mockStorageValue {
		t.Errorf("Expected value %v, got %v", mockStorageValue, result)
	}

	result, ok = ims.Get("nonExistingKey")
	if ok {
		t.Errorf("Expected error 'key not found', got %v", result)

	}

}

func TestSet(t *testing.T) {
	ims := NewStorage()

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
