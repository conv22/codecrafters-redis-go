package reader

import (
	"bufio"
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

type RdbReader struct{}

func (rdb RdbReader) HandleRead(path string) ([]ReaderResult, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	return nil, nil
}
