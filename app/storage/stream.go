package storage

import (
	"sync"
)

type Stream struct {
	mu      sync.Mutex
	entries []*StreamEntry
}

func (s *Stream) AddEntry(entry *StreamEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = append(s.entries, entry)
}

func (s *Stream) GetEntries() []*StreamEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.entries
}

func NewStream() *Stream {
	return &Stream{
		entries: []*StreamEntry{},
	}
}

type StreamEntry struct {
	ID        string
	MsTime    int64
	SqNumber  int64
	KeyValues map[string]string
	mu        sync.Mutex
}

func (e *StreamEntry) AddKeyValuePair(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.KeyValues[key] = value
}

func NewStreamEntry(msTime, sqNumber int64) *StreamEntry {
	return &StreamEntry{
		MsTime:    msTime,
		SqNumber:  sqNumber,
		KeyValues: map[string]string{},
	}
}
