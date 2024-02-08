package parser

type parser_interface interface {
	HandleParse(s string) ([]string, error)
}

func CreateParser(t string) parser_interface {
	switch t {
	case "resp":
		return &RespParser{}

	default:
		return &RespParser{}
	}
}
