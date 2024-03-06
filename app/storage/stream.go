package storage

import (
	"sort"
	"sync"
)

type Stream struct {
	mu         sync.Mutex
	entries    []*StreamEntry
	entriesIds []string
}

func (s *Stream) AddEntry(entry *StreamEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entriesIds = append(s.entriesIds, entry.ID)
	s.entries = append(s.entries, entry)
}

func (s *Stream) GetEntries() []*StreamEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.entries
}

func (s *Stream) GetRange(start, end string) []*StreamEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	startIndex := 0
	endIndex := len(s.entriesIds) - 1

	if start != "-" {
		startIndex = sort.SearchStrings(s.entriesIds, start)
	}

	if end != "+" {
		endIndex = sort.SearchStrings(s.entriesIds, end)
	}

	result := []*StreamEntry{}

	for i := startIndex; i < endIndex; i++ {
		result = append(result, s.entries[i])
	}

	return result

}

func NewStream() *Stream {
	return &Stream{
		entries:    []*StreamEntry{},
		entriesIds: []string{},
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

func NewStreamEntry(id string, msTime, sqNumber int64) *StreamEntry {
	return &StreamEntry{
		ID:        id,
		MsTime:    msTime,
		SqNumber:  sqNumber,
		KeyValues: map[string]string{},
	}
}
