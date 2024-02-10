package parser

import (
	"strconv"
)

func (parser *RespParser) encodeLengthData(encoding string, s string) string {
	return encoding + strconv.Itoa(len(s)) + RespEncodingConstants.Separator + s + RespEncodingConstants.Separator

}

func (parser *RespParser) encodeData(encoding string, s string) string {
	return encoding + s + RespEncodingConstants.Separator
}

func (parser *RespParser) HandleEncode(encoding string, s string) string {
	switch encoding {
	case RespEncodingConstants.String:
		return parser.encodeData(RespEncodingConstants.String, s)
	case RespEncodingConstants.NullBulkString:
		return RespEncodingConstants.NullBulkString + RespEncodingConstants.Separator
	case RespEncodingConstants.Error:
		return parser.encodeData(RespEncodingConstants.Error, s)
	case RespEncodingConstants.BulkString:
		return parser.encodeLengthData(RespEncodingConstants.BulkString, s)
	default:
		return s
	}
}
