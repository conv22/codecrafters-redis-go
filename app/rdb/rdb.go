package rdb

import (
	"bufio"
	"errors"
	"io"
	"os"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type Rdb struct {
	reader             *bufio.Reader
	version            int
	collection         *storage.StorageCollection
	currDbIdx          uint8
	currItemExpiryTime *time.Time
}

func NewRdb() *Rdb {
	return &Rdb{
		collection: storage.NewStorageCollection(),
	}
}

func (rdb *Rdb) setItemToCurrentDB(key storage.StorageKey, encoding string, value interface{}) error {
	currStorage := rdb.collection.GetStorageById(rdb.currDbIdx)

	err := currStorage.Set(key, &storage.StorageItem{
		Expiry:   rdb.currItemExpiryTime,
		Encoding: encoding,
		Value:    value,
	})
	rdb.resetCurrentItemProps()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *Rdb) resetCurrentItemProps() {
	rdb.currItemExpiryTime = nil
}

func (rdb *Rdb) HandleRead(path string) (*storage.StorageCollection, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	rdb.reader = bufio.NewReader(file)

	err = rdb.parseStart()
	if err != nil {
		return nil, err
	}

	version, err := rdb.parseVersion()
	if err != nil {
		return nil, err
	}
	rdb.version = version

out:
	for {
		opCode, err := rdb.readByte()

		currStorage := rdb.collection.GetStorageById(rdb.currDbIdx)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break out
			}
			return nil, err
		}

		switch opCode {
		case RDB_OPCODE_EOF:
			break out
		case RDB_OPCODE_SELECT_DB:
			nextDbIdx, err := rdb.parseSelectDb()
			if err != nil {
				return nil, err
			}
			rdb.currDbIdx = nextDbIdx
		// TODO: skip item if expired
		case RDB_OPCODE_EXPIRE_TIME:
		case RDB_OPCODE_EXPIRE_TIME_MS:
			var time *time.Time
			var err error

			if opCode == RDB_OPCODE_EXPIRE_TIME {
				time, err = rdb.parseExpiryTimeSec()
			} else if opCode == RDB_OPCODE_EXPIRE_TIME_MS {
				time, err = rdb.parseExpiryTimeMs()
			}

			if err != nil {
				return nil, err
			}

			rdb.currItemExpiryTime = time
		case RDB_OPCODE_RESIZE_DB:
			dbHashTableSize, expiryHashTableSize, err := rdb.parseResizeDb()
			if err != nil {
				return nil, err
			}
			currStorage.HashSize = dbHashTableSize
			currStorage.ExpireHashSize = expiryHashTableSize
		case RDB_OPCODE_AUX:
			keyI, value, err := rdb.parseAux()
			if err != nil {
				return nil, err
			}

			key, ok := keyI.(string)

			if !ok {
				continue
			}

			rdb.setItemToCurrentDB(key, "", value)

		default:
			// not supported
			continue
		}
		if err != nil {
			return nil, err
		}
	}

	return rdb.collection, nil

}
