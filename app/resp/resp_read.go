package resp

import (
	"errors"
	"io"
	"strconv"
)

func (r *RespReader) HandleRead() ([]ParsedCmd, int, error) {
	r.bytesRead = 0
	encoding, err := r.readByteAndIncCounter()

	if err != nil {
		if errors.Is(err, io.EOF) {
			return []ParsedCmd{}, 0, nil
		}
		return nil, 0, err
	}

	encStr := string(encoding)

	if encStr == RESP_ENCODING_CONSTANTS.LENGTH {
		result, err := r.parseArray()

		if err != nil {
			return nil, 0, err
		}

		return result, r.bytesRead, nil
	}

	result, err := r.parseByEncoding(encStr)

	if err != nil {
		return nil, 0, err
	}

	return []ParsedCmd{result}, r.bytesRead, nil
}

func (r *RespReader) HandleReadRdbFile() ([]byte, error) {
	r.bytesRead = 0
	_, err := r.readByteAndIncCounter()

	if err != nil {
		return nil, err
	}

	fileLength, err := r.parseInteger()

	if err != nil {
		return nil, err
	}

	return r.readValueByNOfBytesAndIncCounter(fileLength)
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
		encoding, err := r.readByteAndIncCounter()

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
		return ParsedCmd{}, err
	}

	value, err := r.readValueByNOfBytesAndIncCounter(length)

	if err != nil {
		return ParsedCmd{}, err
	}

	if err := r.skipSeparatorAndIncCounter(); err != nil {
		return ParsedCmd{}, err
	}

	return ParsedCmd{ValueType: RESP_ENCODING_CONSTANTS.BULK_STRING, Value: string(value)}, nil
}

func (r *RespReader) parsePrimitiveData(enc string) (ParsedCmd, error) {
	line, err := r.readLineUntilSeperator()

	if err != nil {
		return ParsedCmd{}, err
	}

	return ParsedCmd{ValueType: enc, Value: string(line)}, nil
}

func (r *RespReader) parseInteger() (int, error) {
	line, err := r.readLineUntilSeperator()

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

func (r *RespReader) skipSeparatorAndIncCounter() error {
	bytesSkip, err := r.reader.Discard(len((RESP_ENCODING_CONSTANTS.SEPARATOR)))

	if err != nil {
		return err
	}

	r.bytesRead += bytesSkip

	return nil
}

func (r *RespReader) readValueByNOfBytesAndIncCounter(bytesToRead int) ([]byte, error) {
	data, err := r.reader.Peek(bytesToRead)
	if err != nil {
		return nil, err
	}

	bytesSkip, err := r.reader.Discard(bytesToRead)

	if err != nil {
		return nil, err
	}

	r.bytesRead += bytesSkip

	return data, nil
}

func (r *RespReader) readLineUntilSeperator() (line []byte, err error) {
	line = make([]byte, 0)
	n := 0
	for {
		b, err := r.readByteAndIncCounter()
		if err != nil {
			return nil, err
		}
		n += 1
		line = append(line, b)
		if len(line) > 1 && string(line[len(line)-2:]) == RESP_ENCODING_CONSTANTS.SEPARATOR {
			break
		}
	}

	return line[: len(line)-2 : n], nil
}

func (r *RespReader) readByteAndIncCounter() (byte, error) {
	b, err := r.reader.ReadByte()
	if err != nil {
		return 0, err
	}
	r.bytesRead += 1
	return b, nil

}
