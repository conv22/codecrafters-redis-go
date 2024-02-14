package resp

type RespEncodings struct {
	BulkString     string
	NullBulkString string
	String         string
	Integer        string
	Separator      string
	Length         string
	Error          string
}

var RespEncodingConstants = RespEncodings{
	BulkString:     "$",
	NullBulkString: "$-1",
	String:         "+",
	Integer:        ":",
	Separator:      "\r\n",
	Length:         "*",
	Error:          "-",
}
