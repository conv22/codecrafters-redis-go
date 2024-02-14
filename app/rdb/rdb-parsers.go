package rdb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

// https://rdb.fnordig.de/file_format.html#length-encoding
func (rdb *Rdb) parseLength() (length int, isEncoded bool, err error) {
	firstByte, err := rdb.readByte()
	if err != nil {
		return 0, false, err
	}

	var mask byte = 0b00111111

	switch int((firstByte >> 6)) {
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

		result, err := rdb.readSignedInt()

		if err != nil {
			return 0, false, err
		}

		return int(result), false, nil
	// 11
	case 3:
		return int(firstByte & mask), true, nil

	default:
		return 0, false, nil
	}

}

func (rdb *Rdb) parseStart() error {
	result, err := rdb.readBytes(len(RDB_MAGIC))
	if err != nil {
		return err
	}

	if !bytes.Equal([]byte(RDB_MAGIC), result) {
		return errors.New("redis signature not found")
	}

	return nil
}

func (rdb *Rdb) parseVersion() (int, error) {
	b, err := rdb.readBytes(RDB_VERSION_NUMBER_LENGTH)
	if err != nil {
		return 0, err
	}

	return int(binary.LittleEndian.Uint32(b)), nil
}

// https://github.com/sripathikrishnan/redis-rdb-tools/blob/master/rdbtools/parser.py#L28 extraction algo.
func (rdb *Rdb) handleLZFDecompress(compressed []byte, expectedLength int) ([]byte, error) {
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
func (rdb *Rdb) parseString() (interface{}, error) {
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
			char, err := rdb.readByte()

			if err != nil {
				return "", err
			}
			return char, nil

		// indicates that a 16 bit integer follows
		case 1:
			result, err := rdb.readUnsignedShort()
			if err != nil {
				return "", err
			}
			return result, nil
		// indicates that a 32 bit integer follows
		case 2:
			result, err := rdb.readUnsignedInt()
			if err != nil {
				return "", err
			}
			return result, nil
		// *Compressed Strings
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

			return value, nil

		default:
			return "", errors.New("unsupported encoding")
		}

	} else {
		result, err := rdb.readBytes(length)

		if err != nil {
			return "", err
		}

		return string(result), nil
	}

}

func (rdb *Rdb) parseAux() (key, value interface{}, err error) {
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

func (rdb *Rdb) parseSelectDb() (uint8, error) {
	_, err := rdb.readByte()

	if err != nil {
		return 0, err
	}

	dbNumber, _, err := rdb.parseLength()

	if err != nil {
		return 0, err
	}

	return uint8(dbNumber), nil
}

// The following 4 bytes represent the Unix timestamp as an unsigned integer.
func (rdb *Rdb) parseExpiryTimeSec() (*time.Time, error) {
	result, err := rdb.readBytes(RDB_EXPIRE_TIME_SEC_BYTES_LENGTH)
	if err != nil {
		return nil, err
	}
	unixTimestamp := binary.LittleEndian.Uint32(result)
	expiryTime := time.Unix(int64(unixTimestamp), 0)

	return &expiryTime, nil
}

// The following expiry value is specified in milliseconds. The following 8 bytes represent the Unix timestamp as an unsigned long.
func (rdb *Rdb) parseExpiryTimeMs() (*time.Time, error) {
	result, err := rdb.readBytes(RDB_EXPIRE_TIME_MS_BYTES_LENGTH)
	if err != nil {
		return nil, err
	}
	// Convert the Unix timestamp in milliseconds to a time.Time value
	unixTimestampMs := binary.LittleEndian.Uint64(result)
	expiryTime := time.Unix(0, int64(unixTimestampMs)*int64(time.Millisecond))

	return &expiryTime, nil
}
func (rdb *Rdb) parseResizeDb() (dbHashTableSize, expiryHashTableSize int, err error) {
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
