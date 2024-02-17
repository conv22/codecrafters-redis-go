package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	INFO_CMD_REPLICATION = "replication"
)

func (processor *RespCmdProcessor) handleInfo(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 1 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	switch parsedResult[0].Value {
	case INFO_CMD_REPLICATION:
		replication := processor.replication
		data := []resp.SliceEncoding{
			{S: "role:" + replication.Role, Encoding: resp.RESP_ENCODING_CONSTANTS.SEPARATOR},
			{S: "master_replid:" + replication.MasterReplId, Encoding: resp.RESP_ENCODING_CONSTANTS.SEPARATOR},
			{S: "master_repl_offset:" + replication.Offset, Encoding: resp.RESP_ENCODING_CONSTANTS.SEPARATOR},
		}

		return processor.parser.HandleEncode(RespEncodingConstants.BULK_STRING, processor.parser.HandleEncodeSlices(data))
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "invalid argument")
	}
}
