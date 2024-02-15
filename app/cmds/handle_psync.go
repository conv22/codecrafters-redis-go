package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func (processor *RespCmdProcessor) handlePsync(parsedResult []resp.ParsedCmd) string {
	// TODO
	return processor.parser.HandleEncode(RespEncodingConstants.STRING, "OK")
}
