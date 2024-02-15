package cmds

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (processor *RespCmdProcessor) handlePsync(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 2 {
		processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	builder := strings.Builder{}
	builder.WriteString(CMD_FULL_RESYNC)
	builder.WriteString(" ")
	builder.WriteString(processor.config.Replication.MasterReplId)
	builder.WriteString(" ")
	builder.WriteString(processor.config.Replication.Offset)

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, builder.String())
}
