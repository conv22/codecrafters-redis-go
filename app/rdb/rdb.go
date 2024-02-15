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
	currItemExpiryTime *time.Time
	version            int
	collection         *storage.StorageCollection
}

func NewRdb() *Rdb {
	return &Rdb{
		collection: storage.NewStorageCollection(),
	}
}

func (rdb *Rdb) readObject(encoding byte, key string, currStorage *storage.Storage) error {
	if encoding == RDB_ENCODING_STRING_ENCODING {
		value, err := rdb.parseString()

		if err != nil {
			return err
		}

		currStorage.Set(key, &storage.StorageItem{
			Expiry:   rdb.currItemExpiryTime,
			Value:    value,
			Encoding: encoding,
		})

	}

	return errors.New("not supported encoding")
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

		if err != nil {
			if errors.Is(err, io.EOF) {
				break out
			}
			return nil, err
		}

		currStorage := rdb.collection.GetCurrentStorage()

		switch opCode {
		case RDB_OPCODE_EOF:
			break out
		case RDB_OPCODE_SELECT_DB:
			nextDbIdx, err := rdb.parseSelectDb()
			if err != nil {
				return nil, err
			}
			rdb.collection.CurrStorageId = nextDbIdx
			continue out
		case RDB_OPCODE_RESIZE_DB:
			dbHashTableSize, expiryHashTableSize, err := rdb.parseResizeDb()
			if err != nil {
				return nil, err
			}
			currStorage.HashSize = dbHashTableSize
			currStorage.ExpireHashSize = expiryHashTableSize
			continue out
		case RDB_OPCODE_AUX:
			key, value, err := rdb.parseAux()
			if err != nil {
				return nil, err
			}

			rdb.collection.SetAuxField(key, value)
			continue out
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
		}

		if err != nil {
			return nil, err
		}

		key, err := rdb.parseString()

		if err != nil {
			return nil, err
		}

		rdb.readObject(opCode, key.(string), currStorage)

		rdb.currItemExpiryTime = nil
	}

	return rdb.collection, nil

}
