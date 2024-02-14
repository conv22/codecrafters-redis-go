package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type ReadResult = map[DbKey]*Database

type RdbReader struct {
	reader             *bufio.Reader
	version            int
	dbsMap             map[DbKey]*Database
	currDbIdx          uint8
	currItemExpiryTime *time.Time
}

func NewRdbReader() *RdbReader {
	return &RdbReader{
		dbsMap: make(map[DbKey]*Database),
	}
}

func (rdb *RdbReader) getCurrentDB() *Database {
	db, ok := rdb.dbsMap[rdb.currDbIdx]

	if !ok {
		newDb := newDatabase(rdb.currDbIdx)
		rdb.dbsMap[rdb.currDbIdx] = newDb
		return newDb
	}

	return db
}

func (rdb *RdbReader) setItemToCurrentDB(key ItemKey, encoding string, value interface{}) error {
	db, ok := rdb.dbsMap[rdb.currDbIdx]

	fmt.Printf("%v", value)

	if !ok {
		return errors.New("DB does not exist")
	}
	err := db.setToCache(key, &DbItem{
		expiry:   rdb.currItemExpiryTime,
		encoding: encoding,
		Value:    value,
	})

	if err != nil {
		return err
	}
	rdb.resetCurrentItemProps()
	return nil
}

func (rdb *RdbReader) resetCurrentItemProps() {
	rdb.currItemExpiryTime = nil
}

func (rdb *RdbReader) HandleRead(path string) (*ReadResult, error) {
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

		currDb := rdb.getCurrentDB()

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
			currDb.HashSize = dbHashTableSize
			currDb.ExpireHashSize = expiryHashTableSize
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

	return &rdb.dbsMap, nil

}
