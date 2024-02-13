package reader

import (
	"bufio"
	"errors"
	"io"
	"os"
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

type ReadResult = map[DbKey]*Database

type RdbReader struct {
	reader  *bufio.Reader
	version int
}

func NewRdbReader() *RdbReader {
	return &RdbReader{}
}

func (rdb *RdbReader) HandleRead(path string) (*ReadResult, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var currDbIdx uint8
	dbsMap := make(map[DbKey]*Database)
	var expiryTime *time.Time

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

		currDb, isCurrentDbInCache := dbsMap[currDbIdx]

		if !isCurrentDbInCache {
			currDb = &Database{
				CacheMap: make(map[string]*DbItem),
			}
			dbsMap[currDbIdx] = currDb
		}

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
			currDbIdx = nextDbIdx
		case RDB_OPCODE_EXPIRE_TIME:
			time, err := rdb.parseExpiryTimeSec()
			if err != nil {
				return nil, err
			}
			expiryTime = time
		case RDB_OPCODE_EXPIRE_TIME_MS:
			time, err := rdb.parseExpiryTimeMs()
			if err != nil {
				return nil, err
			}
			expiryTime = time
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

			currDb.CacheMap[key] = &DbItem{
				expiry:   expiryTime,
				encoding: "",
				Value:    value,
			}
			expiryTime = nil
		default:
			// not supported
			continue
		}
		if err != nil {
			return nil, err
		}
	}

	return &dbsMap, nil

}
