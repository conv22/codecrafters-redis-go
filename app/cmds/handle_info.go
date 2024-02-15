package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

const (
	INFO_CMD_REPLICATION = "replication"
)

func (processor *RespCmdProcessor) handleInfo(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 1 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	switch parsedResult[0].Value {
	case INFO_CMD_REPLICATION:
		str := "role:" + processor.config.Role
		return processor.parser.HandleEncode(RespEncodingConstants.BULK_STRING, str)
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "invalid argument")
	}
}
