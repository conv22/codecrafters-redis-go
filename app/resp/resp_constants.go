package resp

type RespEncodings struct {
	BULK_STRING      string
	NULL_BULK_STRING string
	STRING           string
	INTEGER          string
	SEPARATOR        string
	LENGTH           string
	ERROR            string
}

var RESP_ENCODING_CONSTANTS = RespEncodings{
	BULK_STRING:      "$",
	NULL_BULK_STRING: "$-1",
	STRING:           "+",
	INTEGER:          ":",
	SEPARATOR:        "\r\n",
	LENGTH:           "*",
	ERROR:            "-",
}
