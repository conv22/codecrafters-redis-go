package rdb

import (
	"errors"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

var ErrRdbUnsupportedError = errors.New("not supported encoding")

func (rdb *Rdb) readObject(encoding byte, key string, currStorage *storage.Storage) error {
	if encoding == RDB_ENCODING_STRING_ENCODING {
		value, err := rdb.parseString()

		if err != nil {
			return err
		}

		currStorage.Set(key, &storage.StorageItem{
			ExpiryMs: rdb.currItemExpiryTime,
			Value:    value,
			Encoding: encoding,
			Type:     storage.STRING_TYPE,
		})

	}
	// TODO: add other encodings

	return errors.New("not supported encoding")
}
