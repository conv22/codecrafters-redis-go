package resp

import (
	"errors"
	"io"
	"strconv"
)

func (r *RespReader) HandleRead() ([]ParsedCmd, error) {
	defer func() {
		r.result = make([]ParsedCmd, 0)
	}()
	encoding, err := r.reader.ReadByte()

	if err != nil {
		if errors.Is(err, io.EOF) {
			return []ParsedCmd{}, nil
		}
		return nil, err
	}

	encStr := string(encoding)

	if encStr == RESP_ENCODING_CONSTANTS.LENGTH {
		if err := r.parseArray(); err != nil {
			return nil, err
		}
	} else {
		if err := r.parseByEncoding(encStr); err != nil {
			return nil, err
		}
	}

	return r.result, err
}

func (r *RespReader) parseByEncoding(enc string) (err error) {
	switch enc {
	case RESP_ENCODING_CONSTANTS.LENGTH:
		err = r.parseArray()
	case RESP_ENCODING_CONSTANTS.BULK_STRING:
		err = r.parseBulkString()
	case RESP_ENCODING_CONSTANTS.ERROR,
		RESP_ENCODING_CONSTANTS.STRING,
		RESP_ENCODING_CONSTANTS.INTEGER:
		err = r.parsePrimitiveData(enc)
	default:
		err = errors.New("invalid input")
	}

	return err
}

func (r *RespReader) parseArray() error {
	length, err := r.parseInteger()

	if err != nil {
		return err
	}

	for i := 0; i < length; i++ {
		encoding, err := r.reader.ReadByte()

		if err != nil {
			return err
		}

		if err := r.parseByEncoding(string(encoding)); err != nil {
			return err
		}
	}

	return nil
}
func (r *RespReader) parseBulkString() error {
	length, err := r.parseInteger()

	if err != nil {
		return err
	}

	value, err := r.readValueByNOfBytes(length)

	if err != nil {
		return err
	}

	r.result = append(r.result, ParsedCmd{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: string(value)})

	return nil
}

func (r *RespReader) parsePrimitiveData(enc string) error {
	line, _, err := r.readLineUntilSeperator()

	if err != nil {
		return err
	}

	r.result = append(r.result, ParsedCmd{ValueType: enc, Value: string(line)})

	return nil
}

func (r *RespReader) parseInteger() (int, error) {
	line, _, err := r.readLineUntilSeperator()

	if err != nil {
		return 0, err
	}

	strValue := string(line)

	intValue, err := strconv.Atoi(strValue)

	if err != nil {
		return 0, errors.New("invalid integer")
	}

	return intValue, nil
}

func (r *RespReader) skipSeparator() error {
	if _, err := r.reader.Discard(len([]byte(RESP_ENCODING_CONSTANTS.SEPARATOR))); err != nil {
		return err
	}

	return nil
}

func (r *RespReader) readValueByNOfBytes(bytesToRead int) ([]byte, error) {
	data, err := r.reader.Peek(bytesToRead)
	if err != nil {
		return nil, err
	}

	if _, err := r.reader.Discard(bytesToRead); err != nil {
		return nil, err
	}

	if err := r.skipSeparator(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *RespReader) readLineUntilSeperator() (line []byte, bytesRead int, err error) {
	line = make([]byte, 0)
	n := 0
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) > 1 && string(line[len(line)-2:]) == RESP_ENCODING_CONSTANTS.SEPARATOR {
			break
		}
	}
	return line[: len(line)-2 : n], n, nil
}
