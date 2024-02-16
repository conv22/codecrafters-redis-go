package cmds

import (
	"encoding/hex"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const EMPTY_DB_HEX string = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

func (processor *RespCmdProcessor) handlePsync(parsedResult []resp.ParsedCmd) []string {
	if len(parsedResult) < 2 {
		processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	builder := strings.Builder{}
	builder.WriteString(CMD_FULL_RESYNC)
	builder.WriteString(" ")
	builder.WriteString(processor.config.Replication.MasterReplId)
	builder.WriteString(" ")
	builder.WriteString(processor.config.Replication.Offset)

	decoded, err := hex.DecodeString(EMPTY_DB_HEX)

	if err != nil {
		return nil
	}

	ackCmd := processor.parser.HandleEncode(RespEncodingConstants.STRING, builder.String())
	encodingCmd := processor.parser.HandleEncode(RespEncodingConstants.BULK_STRING, string(decoded))

	// exception
	encodingCmd = strings.TrimSuffix(encodingCmd, RespEncodingConstants.SEPARATOR)

	return []string{ackCmd, encodingCmd}
}
