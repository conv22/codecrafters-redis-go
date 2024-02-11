package cmds

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

var RespEncodingConstants = parsers.RespEncodingConstants

type RespCmdProcessor struct {
	parser  parsers.Parser
	storage storage.Storage
	config  config.Config
}

func (processor RespCmdProcessor) ProcessCmd(line string) string {
	parsedResult, err := processor.parser.HandleParse(line)

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "error parsing the line")
	}

	if len(parsedResult) == 0 {
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}

	firstCmd := strings.ToLower(parsedResult[0].Value)
	cmds := parsedResult[1:]

	switch firstCmd {
	case "ping":
		return processor.handlePing()

	case "echo":
		return processor.handleEcho(cmds)

	case "set":
		return processor.handleSet(cmds)

	case "get":
		return processor.handleGet(cmds)

	case "config":
		return processor.handleConfig(cmds)
	case "keys":
		return processor.handleKeys(cmds)
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "not able to process the cmd")
	}
}
