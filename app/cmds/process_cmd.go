package cmds

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

var RespEncodingConstants = resp.RESP_ENCODING_CONSTANTS

type RespCmdProcessor struct {
	parser  *resp.RespParser
	storage *storage.StorageCollection
	config  *config.Config
}

func NewRespCmdProcessor(p *resp.RespParser, storage *storage.StorageCollection, config *config.Config) *RespCmdProcessor {
	return &RespCmdProcessor{
		parser:  p,
		storage: storage,
		config:  config,
	}
}

const (
	CMD_PING   string = "ping"
	CMD_ECHO   string = "echo"
	CMD_GET    string = "get"
	CMD_SET    string = "set"
	CMD_CONFIG string = "config"
	CMD_KEYS   string = "keys"
)

func (processor *RespCmdProcessor) ProcessCmd(line string) string {
	parsedResult, err := processor.parser.HandleParse(line)

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "error parsing the line")
	}

	if len(parsedResult) == 0 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	firstCmd := strings.ToLower(parsedResult[0].Value)
	cmds := parsedResult[1:]

	switch firstCmd {
	case CMD_PING:
		return processor.handlePing()
	case CMD_ECHO:
		return processor.handleEcho(cmds)
	case CMD_SET:
		return processor.handleSet(cmds)
	case CMD_GET:
		return processor.handleGet(cmds)
	case CMD_CONFIG:
		return processor.handleConfig(cmds)
	case CMD_KEYS:
		return processor.handleKeys(cmds)
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not able to process the cmd")
	}
}
