package resp

import (
	"errors"
	"strconv"
	"strings"
)

func parseLength(s string, arrayLength *int, result *[][]ParsedCmd, currIndex int) (string, error) {
	separatorIndex := strings.Index(s, RESP_ENCODING_CONSTANTS.SEPARATOR)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	totalLengthStr := s[0:separatorIndex]
	totalLength, err := strconv.Atoi(totalLengthStr)
	if err != nil {
		return "", errors.New("invalid length")
	}

	*arrayLength = totalLength
	value := s[separatorIndex+len(RESP_ENCODING_CONSTANTS.SEPARATOR):]

	return value, nil
}

func parseLengthData(s string, encoding string, result *[][]ParsedCmd, currIndex int) (string, error) {
	separatorIndex := strings.Index(s, RESP_ENCODING_CONSTANTS.SEPARATOR)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	countStr := s[0:separatorIndex]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", errors.New("invalid length")
	}

	value := s[separatorIndex+len(RESP_ENCODING_CONSTANTS.SEPARATOR):]
	if (len(value) + len(RESP_ENCODING_CONSTANTS.SEPARATOR)) < count {
		return "", errors.New("data length exceeds available data")
	}

	data := value[0:count]
	item := &ParsedCmd{
		Value:     data,
		ValueType: encoding,
	}
	appendItemToResult(result, item, currIndex)

	value = value[count+len(RESP_ENCODING_CONSTANTS.SEPARATOR):]

	return value, nil
}

func parseData(s string, encoding string, result *[][]ParsedCmd, currIndex int) (string, error) {
	separatorIndex := strings.Index(s, RESP_ENCODING_CONSTANTS.SEPARATOR)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	data := s[0:separatorIndex]

	item := &ParsedCmd{
		Value:     data,
		ValueType: encoding,
	}
	value := s[separatorIndex+len(RESP_ENCODING_CONSTANTS.SEPARATOR):]

	appendItemToResult(result, item, currIndex)

	return value, nil
}

func appendItemToResult(result *[][]ParsedCmd, item *ParsedCmd, currIndex int) {
	if len(*result) > currIndex {
		(*result)[currIndex] = append((*result)[currIndex], *item)

	} else {
		newSlice := []ParsedCmd{*item}
		*result = append(*result, newSlice)
	}
}

func parseInteger(s string, result *[][]ParsedCmd, currIndex int) (string, error) {
	return parseData(s, RESP_ENCODING_CONSTANTS.INTEGER, result, currIndex)
}

func parseBulkString(s string, result *[][]ParsedCmd, currIndex int) (string, error) {
	return parseLengthData(s, RESP_ENCODING_CONSTANTS.BULK_STRING, result, currIndex)
}

func parseString(s string, result *[][]ParsedCmd, currIndex int) (string, error) {
	return parseData(s, RESP_ENCODING_CONSTANTS.STRING, result, currIndex)
}

func parseError(s string, result *[][]ParsedCmd, currIndex int) (string, error) {
	return parseData(s, RESP_ENCODING_CONSTANTS.ERROR, result, currIndex)
}

func (parser *RespParser) parseValue(s string, result *[][]ParsedCmd) (string, error) {
	var currArrLength int
	var currIndex int
	firstChar := string(s[0])
	str := s[1:]
	switch firstChar {
	case RESP_ENCODING_CONSTANTS.LENGTH:
		if currArrLength > 0 {
			currArrLength = 0
			currIndex++
		}
		return parseLength(str, &currArrLength, result, currIndex)
	case RESP_ENCODING_CONSTANTS.BULK_STRING:
		return parseBulkString(str, result, currIndex)
	case RESP_ENCODING_CONSTANTS.ERROR:
		return parseError(str, result, currIndex)
	case RESP_ENCODING_CONSTANTS.STRING:
		return parseString(str, result, currIndex)
	case RESP_ENCODING_CONSTANTS.INTEGER:
		return parseInteger(str, result, currIndex)
	default:
		return "", errors.New("invalid input")
	}
}

func (parser *RespParser) HandleParse(s string) ([][]ParsedCmd, error) {
	var result = [][]ParsedCmd{}

	if len(s) == 0 {
		return result, nil
	}
	if !strings.HasSuffix(s, RESP_ENCODING_CONSTANTS.SEPARATOR) {
		return nil, errors.New("invalid input")
	}

	str := strings.Clone(s)

	for len(str) != 0 {
		parsedValue, err := parser.parseValue(str, &result)
		if err != nil {
			return nil, err
		}
		str = parsedValue

	}

	return result, nil

}
