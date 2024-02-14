package resp

type ParsedCmd struct {
	ValueType string
	Value     string
}

type RespParser struct{}

func NewRespParser() *RespParser {
	return &RespParser{}
}
