package resp

import (
	"strconv"
	"strings"
)

func encodeLengthData(encoding string, s string) string {
	return encoding + strconv.Itoa(len([]byte(s))) + RESP_ENCODING_CONSTANTS.SEPARATOR + s + RESP_ENCODING_CONSTANTS.SEPARATOR
}

func encodeData(encoding string, s string) string {
	return encoding + s + RESP_ENCODING_CONSTANTS.SEPARATOR
}

type SliceEncoding struct {
	S        string
	Encoding string
}

func (parser RespParser) HandleEncodeSliceList(slices []SliceEncoding) string {
	length := strconv.Itoa(len(slices))
	return parser.handleEncodeSlices(slices, RESP_ENCODING_CONSTANTS.LENGTH+length+RESP_ENCODING_CONSTANTS.SEPARATOR)
}

func (parser RespParser) HandleEncodeSlices(slices []SliceEncoding) string {
	return parser.handleEncodeSlices(slices, "")
}
func (parser RespParser) handleEncodeSlices(slices []SliceEncoding, prefix string) string {
	builder := strings.Builder{}
	if prefix != "" {
		builder.WriteString(prefix)
	}

	for _, slice := range slices {
		builder.WriteString(parser.HandleEncode(slice.Encoding, slice.S))
	}

	return builder.String()
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
	case RESP_ENCODING_CONSTANTS.SEPARATOR:
		return s + RESP_ENCODING_CONSTANTS.SEPARATOR
	default:
		return ""
	}
}
