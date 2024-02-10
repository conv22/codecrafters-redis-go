package parser

import (
	"errors"
	"strconv"
	"strings"
)

func (parser *RespParser) parseLength(s string, arrayLength *int, result *[]ParsedCmd) (string, error) {
	if *arrayLength > 0 {
		return "", errors.New("invalid input")
	}

	separatorIndex := strings.Index(s, RespEncodingConstants.Separator)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	totalLengthStr := s[0:separatorIndex]
	totalLength, err := strconv.Atoi(totalLengthStr)
	if err != nil {
		return "", errors.New("invalid length")
	}

	*arrayLength = totalLength
	value := s[separatorIndex+len(RespEncodingConstants.Separator):]

	return value, nil
}

func (parser *RespParser) parseLengthData(s string, encoding string, result *[]ParsedCmd) (string, error) {
	separatorIndex := strings.Index(s, RespEncodingConstants.Separator)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	countStr := s[0:separatorIndex]
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return "", errors.New("invalid length")
	}

	value := s[separatorIndex+len(RespEncodingConstants.Separator):]
	if (len(value) + len(RespEncodingConstants.Separator)) < count {
		return "", errors.New("data length exceeds available data")
	}

	data := value[0:count]
	item := &ParsedCmd{
		Value:     data,
		ValueType: encoding,
	}
	*result = append(*result, *item)
	value = value[count+len(RespEncodingConstants.Separator):]

	return value, nil
}

func (parser *RespParser) parseData(s string, encoding string, result *[]ParsedCmd) (string, error) {
	separatorIndex := strings.Index(s, RespEncodingConstants.Separator)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	data := s[0:separatorIndex]

	item := &ParsedCmd{
		Value:     data,
		ValueType: encoding,
	}
	value := s[separatorIndex+len(RespEncodingConstants.Separator):]

	*result = append(*result, *item)

	return value, nil
}

func (parser *RespParser) parseInteger(s string, result *[]ParsedCmd) (string, error) {
	return parser.parseData(s, RespEncodingConstants.Integer, result)
}

func (parser *RespParser) parseBulkString(s string, result *[]ParsedCmd) (string, error) {
	return parser.parseLengthData(s, RespEncodingConstants.BulkString, result)
}

func (parser *RespParser) parseString(s string, result *[]ParsedCmd) (string, error) {
	return parser.parseData(s, RespEncodingConstants.String, result)
}

func (parser *RespParser) parseError(s string, result *[]ParsedCmd) (string, error) {
	return parser.parseData(s, RespEncodingConstants.Error, result)
}

func (parser *RespParser) parseValue(s string, arrayLength *int, result *[]ParsedCmd) (string, error) {

	firstChar := string(s[0])
	str := s[1:]
	switch firstChar {
	case RespEncodingConstants.Length:
		return parser.parseLength(str, arrayLength, result)
	case RespEncodingConstants.BulkString:
		return parser.parseBulkString(str, result)
	case RespEncodingConstants.Error:
		return parser.parseError(str, result)
	case RespEncodingConstants.String:
		return parser.parseString(str, result)
	case RespEncodingConstants.Integer:
		return parser.parseInteger(str, result)
	default:
		return "", errors.New("invalid input")
	}
}

func (parser *RespParser) HandleParse(s string) ([]ParsedCmd, error) {
	var arrayLength int
	var result = []ParsedCmd{}

	if len(s) == 0 {
		return result, nil
	}
	if !strings.HasSuffix(s, RespEncodingConstants.Separator) {
		return nil, errors.New("invalid input")
	}

	str := strings.Clone(s)

	for len(str) != 0 {
		parsedValue, err := parser.parseValue(str, &arrayLength, &result)
		if err != nil {
			return nil, err
		}
		str = parsedValue

	}

	if arrayLength != 0 && arrayLength != len(result) {
		return nil, errors.New("invalid input")
	}

	return result, nil

}
