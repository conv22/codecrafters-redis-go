package parser

type parser_interface interface {
	HandleParse(s string) ([]string, error)
	HandleEncode(encoding string, s string) string
}

func CreateParser(t string) parser_interface {
	switch t {
	case "resp":
		return &RespParser{}

	default:
		return &RespParser{}
	}
}
