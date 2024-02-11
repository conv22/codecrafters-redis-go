package parser

import (
	"strconv"
)

func encodeLengthData(encoding string, s string) string {
	return encoding + strconv.Itoa(len(s)) + RespEncodingConstants.Separator + s + RespEncodingConstants.Separator
}

func encodeData(encoding string, s string) string {
	return encoding + s + RespEncodingConstants.Separator
}

type SliceEncoding struct {
	S        string
	Encoding string
}

func (parser RespParser) HandleEncodeSlice(slices []SliceEncoding) string {
	length := strconv.Itoa(len(slices))
	output := RespEncodingConstants.Length + length + RespEncodingConstants.Separator

	for _, slice := range slices {
		encodedValue := parser.HandleEncode(slice.Encoding, slice.S)
		output += encodedValue
	}

	return output
}

func (parser RespParser) HandleEncode(encoding string, s string) string {
	switch encoding {
	case RespEncodingConstants.String:
		return encodeData(RespEncodingConstants.String, s)
	case RespEncodingConstants.NullBulkString:
		return RespEncodingConstants.NullBulkString + RespEncodingConstants.Separator
	case RespEncodingConstants.Error:
		return encodeData(RespEncodingConstants.Error, s)
	case RespEncodingConstants.BulkString:
		return encodeLengthData(RespEncodingConstants.BulkString, s)
	default:
		return ""
	}
}
