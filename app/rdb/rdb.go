package rdb

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type Rdb struct {
	reader             *bufio.Reader
	currItemExpiryTime int64
	version            int
	collection         *storage.StorageCollection
}

func NewRdb() *Rdb {
	return &Rdb{
		collection: storage.NewStorageCollection(),
	}
}

func (rdb *Rdb) HandleReadFromBytes(data []byte) (*storage.StorageCollection, error) {
	reader := bytes.NewReader(data)

	rdb.reader = bufio.NewReader(reader)

	return rdb.handleRead()
}

func (rdb *Rdb) HandleReadFromFile(path string) (*storage.StorageCollection, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	rdb.reader = bufio.NewReader(file)

	return rdb.handleRead()
}

func (rdb *Rdb) handleRead() (*storage.StorageCollection, error) {
	err := rdb.parseStart()
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
		case RDB_OPCODE_EXPIRE_TIME, RDB_OPCODE_EXPIRE_TIME_MS:
			var time int64
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
			continue out
		}

		if err != nil {
			return nil, err
		}

		parseStr, err := rdb.parseString()

		if err != nil {
			return nil, err
		}

		key, ok := parseStr.(string)

		if !ok {
			return nil, errors.New("invalid encoding")
		}

		err = rdb.readObject(opCode, key, currStorage)

		if errors.Is(ErrRdbUnsupportedError, err) {
			rdb.skipObject(opCode)
		}

		rdb.currItemExpiryTime = 0
	}

	return rdb.collection, nil

}
