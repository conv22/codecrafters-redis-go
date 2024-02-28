package resp

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

func (r *RespReader) HandleRead() ([]ParsedCmd, error) {
	encoding, err := r.reader.ReadByte()

	if err != nil {
		if errors.Is(err, io.EOF) {
			return []ParsedCmd{}, nil
		}
		return nil, err
	}

	encStr := string(encoding)

	if encStr == RESP_ENCODING_CONSTANTS.LENGTH {
		result, err := r.parseArray()

		if err != nil {
			return nil, err
		}

		return result, nil
	}

	result, err := r.parseByEncoding(encStr)

	if err != nil {
		return nil, err
	}

	return []ParsedCmd{result}, nil

}

func (r *RespReader) HandleReadRdbFile() ([]byte, error) {
	_, err := r.reader.ReadByte()

	if err != nil {
		return nil, err
	}

	fileLength, err := r.parseInteger()

	if err != nil {
		return nil, err
	}

	return r.readValueByNOfBytes(fileLength)
}

func (r *RespReader) parseByEncoding(enc string) (result ParsedCmd, err error) {
	switch enc {
	case RESP_ENCODING_CONSTANTS.BULK_STRING:
		result, err = r.parseBulkString()
	case RESP_ENCODING_CONSTANTS.ERROR,
		RESP_ENCODING_CONSTANTS.STRING,
		RESP_ENCODING_CONSTANTS.INTEGER:
		result, err = r.parsePrimitiveData(enc)
	default:
		err = errors.New("invalid input")
	}

	return result, err
}

func (r *RespReader) parseArray() (result []ParsedCmd, err error) {
	length, err := r.parseInteger()

	if err != nil {
		return nil, err
	}

	result = []ParsedCmd{}

	for i := 0; i < length; i++ {
		encoding, err := r.reader.ReadByte()

		if err != nil {
			return nil, err
		}

		parsed, err := r.parseByEncoding(string(encoding))

		if err != nil {
			return nil, err
		}

		result = append(result, parsed)
	}

	return result, err
}
func (r *RespReader) parseBulkString() (ParsedCmd, error) {
	length, err := r.parseInteger()

	if err != nil {
		fmt.Println("INT ERROR", err)
		return ParsedCmd{}, err
	}

	value, err := r.readValueByNOfBytes(length)

	if err != nil {
		fmt.Println("readValueByNOfBytes ERROR", err)

		return ParsedCmd{}, err
	}

	if err := r.skipSeparator(); err != nil {
		fmt.Println("skipSeparator ERROR", err)

		return ParsedCmd{}, err
	}

	return ParsedCmd{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: string(value)}, nil
}

func (r *RespReader) parsePrimitiveData(enc string) (ParsedCmd, error) {
	line, _, err := r.readLineUntilSeperator()

	if err != nil {
		return ParsedCmd{}, err
	}

	return ParsedCmd{ValueType: enc, Value: string(line)}, nil
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
