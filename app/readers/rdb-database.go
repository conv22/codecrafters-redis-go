package reader

import (
	"time"
)

type ItemKey = string

type DbKey = uint8

type DbItem struct {
	Value    any
	expiry   *time.Time
	encoding string
}

type Database struct {
	ID             uint8
	HashSize       int
	ExpireHashSize int
	CacheMap       map[ItemKey]*DbItem
}

func newDatabase(ID uint8) *Database {
	return &Database{
		ID:       ID,
		CacheMap: make(map[string]*DbItem),
	}
}

func (db *Database) setToCache(key ItemKey, item *DbItem) error {
	db.CacheMap[key] = item

	return nil
}
