package parser

type ParsedCmd struct {
	ValueType string
	Value     string
}

type Parser interface {
	HandleParse(s string) ([]ParsedCmd, error)
	HandleEncode(encoding string, s string) string
	HandleEncodeSlice(encoding []SliceEncoding) string
}

func CreateParser(t string) Parser {
	switch t {
	case "resp":
		return &RespParser{}

	default:
		return &RespParser{}
	}
}
