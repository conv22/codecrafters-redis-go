package resp

import (
	"bufio"
	"io"
)

type ParsedCmd struct {
	ValueType string
	Value     string
}

type RespReader struct {
	reader *bufio.Reader
	result []ParsedCmd
}

func NewReader(reader io.Reader) *RespReader {
	return &RespReader{
		reader: bufio.NewReader(reader),
		result: make([]ParsedCmd, 0),
	}
}
