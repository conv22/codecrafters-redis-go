package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func (processor *RespCmdProcessor) handleEcho(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")

	}
	return processor.parser.HandleEncode(RespEncodingConstants.STRING, parsedResult[0].Value)
}
