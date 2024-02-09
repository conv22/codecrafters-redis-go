package parser

import (
	"errors"
	"strconv"
	"strings"
)

type RespEncodings struct {
	BulkString string
	String     string
	Integer    string
	Separator  string
	Length     string
	Error      string
}

var RespEncodingConstants = RespEncodings{
	BulkString: "$",
	String:     "+",
	Integer:    ":",
	Separator:  "\r\n",
	Length:     "*",
	Error:      "-",
}

type RespParser struct {
	totalLength int
	result      []string
}

func (parser *RespParser) parseLength(s string) (string, error) {
	if parser.totalLength > 0 {
		return "", errors.New("invalid input")
	}

	separatorIndex := strings.Index(s, RespEncodingConstants.Separator)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	totalLengthStr := s[1:separatorIndex]
	totalLength, err := strconv.Atoi(totalLengthStr)
	if err != nil {
		return "", errors.New("invalid length")
	}

	parser.totalLength = totalLength
	value := s[separatorIndex+len(RespEncodingConstants.Separator):]

	return value, nil
}

func (parser *RespParser) parseBulkString(s string) (string, error) {
	separatorIndex := strings.Index(s, RespEncodingConstants.Separator)
	if separatorIndex == -1 {
		return "", errors.New("separator not found")
	}

	charCountStr := s[1:separatorIndex]
	charCount, err := strconv.Atoi(charCountStr)
	if err != nil {
		return "", errors.New("invalid length")
	}

	value := s[separatorIndex+len(RespEncodingConstants.Separator):]
	if (len(value) + len(RespEncodingConstants.Separator)) < charCount {
		return "", errors.New("bulk string length exceeds available data")
	}

	word := value[0:charCount]
	parser.result = append(parser.result, word)
	value = value[charCount+len(RespEncodingConstants.Separator):]

	return value, nil
}

func (parser *RespParser) parseValue(s string) (string, error) {
	firstChar := string(s[0])
	switch firstChar {
	case RespEncodingConstants.Length:
		return parser.parseLength(s)
	case RespEncodingConstants.BulkString:
		return parser.parseBulkString(s)
	default:
		return "", errors.New("invalid input")
	}
}

func (parser *RespParser) resetParser() {
	parser.totalLength = 0
	parser.result = []string{}
}

func (parser *RespParser) HandleEncode(encoding string, s string) string {
	switch encoding {
	case RespEncodingConstants.String:
		return RespEncodingConstants.String + s + RespEncodingConstants.Separator
	}
	return s
}

func (parser *RespParser) HandleParse(s string) ([]string, error) {
	parser.resetParser()

	if len(s) == 0 {
		return parser.result, nil
	}
	if !strings.HasPrefix(s, RespEncodingConstants.Length) {
		return nil, errors.New("invalid input")
	}

	str := strings.Clone(s)

	for len(str) != 0 {
		parsedValue, err := parser.parseValue(str)
		if err != nil {
			return nil, err
		}
		str = parsedValue

	}

	if parser.totalLength != len(parser.result) {
		return nil, errors.New("invalid input")
	}

	return parser.result, nil

}
