package storage

import (
	"sync"
)

type storageKey = string

type storageId = uint8

const (
	STRING_TYPE = "string"
	NONE_TYPE   = "none"
	STREAM      = "stream"
)

type StreamEntry struct {
	ID        string
	MsTime    int64
	SqNumber  int64
	KeyValues map[string]interface{}
	mu        sync.Mutex
}

func (e *StreamEntry) AddEntry(key string, value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.KeyValues[key] = value
}

func NewStreamEntry(id string, stream *Stream, msTime, sqNumber int64) *StreamEntry {
	return &StreamEntry{
		ID:        id,
		MsTime:    msTime,
		SqNumber:  sqNumber,
		KeyValues: map[string]interface{}{},
	}
}

type Stream struct {
	mu      sync.Mutex
	Entries []*StreamEntry
}

func (s *Stream) AddEntry(entry *StreamEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Entries = append(s.Entries, entry)
}

func NewStream() *Stream {
	return &Stream{
		Entries: []*StreamEntry{},
	}
}

type StorageItem struct {
	Value    any
	Type     string
	ExpiryMs int64
	Encoding byte
}

type Storage struct {
	HashSize       int
	ExpireHashSize int
	CacheMap       map[storageKey]*StorageItem
	mu             sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		CacheMap: make(map[storageKey]*StorageItem),
	}
}

func (ims *Storage) Get(key storageKey) (*StorageItem, bool) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	value, ok := ims.CacheMap[key]
	return value, ok
}

func (ims *Storage) GetKeys() []string {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	keys := make([]string, 0, len(ims.CacheMap))
	for k := range ims.CacheMap {
		keys = append(keys, k)
	}
	return keys
}

func (ims *Storage) Set(key storageKey, value *StorageItem) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.CacheMap[key] = value
	return nil
}

func (ims *Storage) Delete(key storageKey) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.CacheMap, key)
	return nil
}
