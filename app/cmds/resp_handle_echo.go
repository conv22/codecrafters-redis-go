package cmds

import (
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
)

func (processor *RespCmdProcessor) handleEcho(parsedResult []parsers.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")

	}
	return processor.parser.HandleEncode(RespEncodingConstants.String, parsedResult[0].Value)
}
