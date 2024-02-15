package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func (processor *RespCmdProcessor) handleReplConf(parsedResult []resp.ParsedCmd) string {
	return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
}
