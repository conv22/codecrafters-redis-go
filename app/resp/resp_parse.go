package resp

import (
	"errors"
	"strconv"
	"strings"
)

func parseLength(s string) (nextStr string, length int, err error) {
	separatorIndex := strings.Index(s, RESP_ENCODING_CONSTANTS.SEPARATOR)
	if separatorIndex == -1 {
		return "", 0, errors.New("separator not found")
	}

	totalLengthStr := s[:separatorIndex]
	totalLength, err := strconv.Atoi(totalLengthStr)
	if err != nil {
		return "", 0, errors.New("invalid length")
	}

	return s[separatorIndex+len(RESP_ENCODING_CONSTANTS.SEPARATOR):], totalLength, nil
}

func parseLengthData(s string) (nextStr, value string, err error) {
	separatorIndex := strings.Index(s, RESP_ENCODING_CONSTANTS.SEPARATOR)
	if separatorIndex == -1 {
		return "", "", errors.New("separator not found")
	}

	countStr := s[0:separatorIndex]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", "", errors.New("invalid length")
	}

	value = s[separatorIndex+len(RESP_ENCODING_CONSTANTS.SEPARATOR):]
	if (len(value) + len(RESP_ENCODING_CONSTANTS.SEPARATOR)) < count {
		return "", "", errors.New("data length exceeds available data")
	}

	return value[count+len(RESP_ENCODING_CONSTANTS.SEPARATOR):], value[:count], nil
}

func parseData(s string) (nextStr, value string, err error) {
	separatorIndex := strings.Index(s, RESP_ENCODING_CONSTANTS.SEPARATOR)
	if separatorIndex == -1 {
		return "", "", errors.New("separator not found")
	}

	return s[separatorIndex+len(RESP_ENCODING_CONSTANTS.SEPARATOR):], s[:separatorIndex], nil
}

func parseValue(s string, result *[]ParsedCmd, arrLength *int) (nextStr string, err error) {
	var value string
	encoding := string(s[0])
	str := s[1:]
	switch encoding {
	case RESP_ENCODING_CONSTANTS.LENGTH:
		nextStr, length, err := parseLength(str)
		*arrLength = length
		return nextStr, err
	case RESP_ENCODING_CONSTANTS.BULK_STRING:
		nextStr, value, err = parseLengthData(str)
	case RESP_ENCODING_CONSTANTS.ERROR,
		RESP_ENCODING_CONSTANTS.STRING,
		RESP_ENCODING_CONSTANTS.INTEGER:
		nextStr, value, err = parseData(str)
	default:
		return "", errors.New("invalid input")
	}

	item := ParsedCmd{
		ValueType: encoding,
		Value:     value,
	}

	*result = append(*result, item)

	return nextStr, err
}

func HandleParse(s string) ([]ParsedCmd, error) {
	result := []ParsedCmd{}
	var arrLength int

	if len(s) == 0 {
		return result, nil
	}
	if !strings.HasSuffix(s, RESP_ENCODING_CONSTANTS.SEPARATOR) {
		return nil, errors.New("invalid input")
	}

	str := strings.Clone(s)
	for len(str) != 0 {
		parsedValue, err := parseValue(str, &result, &arrLength)
		if err != nil {
			return nil, err
		}
		str = parsedValue

	}

	if arrLength > 0 && arrLength != len(result) {
		return nil, errors.New("invalid array length specified")
	}

	return result, nil
}
