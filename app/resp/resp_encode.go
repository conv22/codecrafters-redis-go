package resp

import (
	"strconv"
)

func encodeLengthData(encoding string, s string) string {
	return encoding + strconv.Itoa(len(s)) + RESP_ENCODING_CONSTANTS.SEPARATOR + s + RESP_ENCODING_CONSTANTS.SEPARATOR
}

func encodeData(encoding string, s string) string {
	return encoding + s + RESP_ENCODING_CONSTANTS.SEPARATOR
}

type SliceEncoding struct {
	S        string
	Encoding string
}

func (parser RespParser) HandleEncodeSlice(slices []SliceEncoding) string {
	length := strconv.Itoa(len(slices))
	output := RESP_ENCODING_CONSTANTS.LENGTH + length + RESP_ENCODING_CONSTANTS.SEPARATOR

	for _, slice := range slices {
		encodedValue := parser.HandleEncode(slice.Encoding, slice.S)
		output += encodedValue
	}

	return output
}

func (parser RespParser) HandleEncode(encoding string, s string) string {
	switch encoding {
	case RESP_ENCODING_CONSTANTS.STRING:
		return encodeData(RESP_ENCODING_CONSTANTS.STRING, s)
	case RESP_ENCODING_CONSTANTS.NULL_BULK_STRING:
		return RESP_ENCODING_CONSTANTS.NULL_BULK_STRING + RESP_ENCODING_CONSTANTS.SEPARATOR
	case RESP_ENCODING_CONSTANTS.ERROR:
		return encodeData(RESP_ENCODING_CONSTANTS.ERROR, s)
	case RESP_ENCODING_CONSTANTS.BULK_STRING:
		return encodeLengthData(RESP_ENCODING_CONSTANTS.BULK_STRING, s)
	default:
		return ""
	}
}
