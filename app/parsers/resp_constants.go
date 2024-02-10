package parser

type RespEncodings struct {
	BulkString string
	String     string
	Integer    string
	Separator  string
	Length     string
	Error      string
	Null       string
}

var RespEncodingConstants = RespEncodings{
	BulkString: "$",
	String:     "+",
	Integer:    ":",
	Separator:  "\r\n",
	Length:     "*",
	Error:      "-",
	Null:       "_",
}
