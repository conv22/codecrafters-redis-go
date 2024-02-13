package reader

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// https://rdb.fnordig.de/file_format.html#:~:text=0003%22%20=%3E%20Version%203-,Op%20Codes,-Each%20part%20after
const (
	// End of the RDB file
	EOF = 0xFF
	// Database selector
	SELECTDB = 0xFE
	// Expire time in seconds
	EXPIRETIME = 0xFD
	// Expire time in milliseconds
	EXPIRETIMEMS = 0xFC
	// Hash table sizes for main keyspace and expires
	RESIZEDB = 0xFB
	// Auxiliary fields. Arbitrary key-value settings
	AUX = 0xFA
	// String encoded type
	TYPE_STRING = 0
)

const (
	// Length of version in the header
	VERSION_NUMBER_LENGTH = 4
	// Number of byes for EXPIRETIMEMS
	EXPIRE_TIME_MS_BYTES_LENGTH = 8
	// Number of byes for EXPIRETIME
	EXPIRE_TIME_SEC_BYTES_LENGTH = 4
)

// Magic redis string
var RDB_MAGIC = []byte("REDIS")

type RdbReader struct {
	reader         *bufio.Reader
	expirationTime *time.Time
}

func (rdb *RdbReader) readBytes(l int) ([]byte, error) {
	buffer := make([]byte, l)

	length, err := rdb.reader.Read(buffer)

	if err != nil {
		return nil, err
	}

	if length != l {
		return nil, errors.New("invalid encoding")
	}

	return buffer, nil

}

func (rdb *RdbReader) readByte() (byte, error) {
	bytes, err := rdb.readBytes(1)
	if err != nil {
		return 0, err
	}
	return bytes[0], nil

}

// https://rdb.fnordig.de/file_format.html#length-encoding
func (rdb *RdbReader) parseLength() (length int, isEncoded bool, err error) {
	firstByte, err := rdb.readByte()
	if err != nil {
		return 0, false, err
	}

	msb := int(firstByte >> 6)
	var mask byte = 0b00111111

	switch msb {
	// 00
	case 0:
		return int(firstByte & mask), false, nil
	// 01
	case 1:
		secondByte, err := rdb.readByte()
		if err != nil {
			return 0, false, err
		}
		return int((firstByte&mask)<<8) | int(secondByte), false, nil
	// 10
	case 2:

		result, err := rdb.readBytes(4)

		if err != nil {
			return 0, false, err
		}

		return int(binary.LittleEndian.Uint32(result)), false, nil
	// 11
	case 3:
		return int(firstByte & mask), true, nil

	default:
		return 0, false, nil
	}

}

func (rdb *RdbReader) parseStart() error {
	result, err := rdb.readBytes(len(RDB_MAGIC))
	if err != nil {
		return err
	}

	if !bytes.Equal(RDB_MAGIC, result) {
		return errors.New("redis signature not found")
	}

	return nil
}

func (rdb *RdbReader) parseVersion() error {
	_, err := rdb.readBytes(VERSION_NUMBER_LENGTH)
	if err != nil {
		return err
	}

	return nil
}

// https://github.com/sripathikrishnan/redis-rdb-tools/blob/master/rdbtools/parser.py#L28 extraction algo.
func (rdb *RdbReader) handleLZFDecompress(compressed []byte, expectedLength int) ([]byte, error) {
	inLen := len(compressed)
	inIndex := 0
	var outStream []byte
	outIndex := 0

	for inIndex < inLen {
		ctrl := compressed[inIndex]
		inIndex++

		if ctrl < 32 {
			for x := 0; x <= int(ctrl); x++ {
				outStream = append(outStream, compressed[inIndex])
				inIndex++
				outIndex++
			}
		} else {
			length := int(ctrl >> 5)
			if length == 7 {
				length += int(compressed[inIndex])
				inIndex++
			}

			ref := outIndex - ((int(ctrl) & 0x1f) << 8) - int(compressed[inIndex]) - 1
			inIndex++

			for x := 0; x < length+2; x++ {
				outStream = append(outStream, outStream[ref])
				ref++
				outIndex++
			}
		}
	}

	if len(outStream) != expectedLength {
		return nil, errors.New("invalid input")
	}

	return outStream, nil
}
func (rdb *RdbReader) parseString() (string, error) {
	length, isEncoded, err := rdb.parseLength()

	if err != nil {
		return "", err
	}

	if isEncoded {
		// https://rdb.fnordig.de/file_format.html#string-encoding
		switch length {
		// *Integers as String
		// indicates that an 8 bit integer follows
		case 0:
			length = 1
		// indicates that a 16 bit integer follows
		case 1:
			length = 2
		// indicates that a 32 bit integer follows
		case 2:
			length = 4
		// *Compressed Strings, unsupported for now
		case 3:
			clenLength, _, err := rdb.parseLength()
			if err != nil {
				return "", err
			}

			l, _, err := rdb.parseLength()

			if err != nil {
				return "", err
			}

			clenValue, err := rdb.readBytes(clenLength)

			if err != nil {
				return "", err
			}

			value, err := rdb.handleLZFDecompress(clenValue, l)

			if err != nil {
				return "", err
			}

			return string(value), nil

		default:
			return "", errors.New("unsupported encoding")
		}

	}
	result, err := rdb.readBytes(length)

	if err != nil {
		return "", err
	}

	fmt.Println(string(result))

	return string(result), nil

}

func (rdb *RdbReader) parseAux() (key, value string, err error) {
	key, err = rdb.parseString()
	if err != nil {
		return
	}
	value, err = rdb.parseString()
	if err != nil {
		return
	}

	return
}

func (rdb *RdbReader) parseSelectDb() error {
	_, err := rdb.readByte()

	if err != nil {
		return err
	}

	_, _, err = rdb.parseLength()

	if err != nil {
		return err
	}

	return nil
}

// The following 4 bytes represent the Unix timestamp as an unsigned integer.
func (rdb *RdbReader) parseExpiryTimeSec() (uint32, error) {
	result, err := rdb.readBytes(EXPIRE_TIME_SEC_BYTES_LENGTH)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(result), nil
}

// The following expiry value is specified in milliseconds. The following 8 bytes represent the Unix timestamp as an unsigned long.
func (rdb *RdbReader) parseExpiryTimeMs() (uint64, error) {
	result, err := rdb.readBytes(EXPIRE_TIME_MS_BYTES_LENGTH)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint64(result), nil
}
func (rdb *RdbReader) parseResizeDb() (dbHashTableSize, expiryHashTableSize int, err error) {
	dbHashTableSize, _, err = rdb.parseLength()
	if err != nil {
		return
	}

	expiryHashTableSize, _, err = rdb.parseLength()

	if err != nil {
		return
	}

	return
}

func (rdb *RdbReader) HandleRead(path string) ([]ReaderResult, error) {
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

	err = rdb.parseVersion()
	if err != nil {
		return nil, err
	}

	var results []ReaderResult

out:
	for {
		opCode, err := rdb.readByte()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break out
			}
			return nil, err
		}

		switch opCode {

		case EOF:
			break out
		case SELECTDB:
			err = rdb.parseSelectDb()
		case EXPIRETIME:
			timeStamp, err := rdb.parseExpiryTimeSec()
			if err != nil {
				return nil, err
			}
			expiryTime := time.Unix(int64(timeStamp), 0)
			rdb.expirationTime = &expiryTime
		case EXPIRETIMEMS:
			timeStamp, err := rdb.parseExpiryTimeMs()
			if err != nil {
				return nil, err
			}
			expiryTime := time.UnixMilli(int64(timeStamp))
			rdb.expirationTime = &expiryTime
		case RESIZEDB:
			_, _, err = rdb.parseResizeDb()
		case AUX:
			key, value, err := rdb.parseAux()
			if err != nil {
				return nil, err
			}
			results = append(results, ReaderResult{Key: key, Value: value})
		case TYPE_STRING:
			str, err := rdb.parseString()
			if err != nil {
				return nil, err
			}
			results = append(results, ReaderResult{Key: str, Value: ""})

		default:
			continue
		}
		if err != nil {
			return nil, err
		}

	}

	return results, nil

}
