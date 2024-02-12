package reader

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
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

// Length of version in the header
const VERSION_NUMBER_LENGTH = 4

// Magic redis string
var RDB_MAGIC = []byte("REDIS")

type RdbReader struct {
	reader *bufio.Reader
}

func (rdb *RdbReader) parseStart() error {
	buffer := make([]byte, len(RDB_MAGIC))
	length, err := rdb.reader.Read(buffer)
	if err != nil {
		return err
	}

	if length != len(RDB_MAGIC) || bytes.Equal(RDB_MAGIC, buffer) {
		return errors.New("redis signature not found")
	}

	return nil
}
func (rdb *RdbReader) parseVersion() error {
	buffer := make([]byte, VERSION_NUMBER_LENGTH)
	length, err := rdb.reader.Read(buffer)
	if err != nil {
		return err
	}

	if length != VERSION_NUMBER_LENGTH {
		return errors.New("redis signature not found")
	}

	return nil
}

func (rdb RdbReader) HandleRead(path string) ([]ReaderResult, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

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
		opCode, err := reader.ReadByte()

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		switch opCode {
		case EOF:
			break out
		case SELECTDB:
		case EXPIRETIME:
		case EXPIRETIMEMS:
		case RESIZEDB:
		case AUX:
		case TYPE_STRING:
		default:
			break out
		}

	}
	return results, nil

}
