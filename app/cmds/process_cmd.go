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
	CMD_PING        string = "PING"
	CMD_PONG        string = "PONG"
	CMD_ECHO        string = "ECHO"
	CMD_GET         string = "GET"
	CMD_SET         string = "SET"
	CMD_CONFIG      string = "CONFIG"
	CMD_KEYS        string = "KEYS"
	CMD_INFO        string = "INFO"
	CMD_REPLCONF    string = "REPLCONF"
	CMD_PSYNC       string = "PSYNC"
	CMD_OK          string = "OK"
	CMD_FULL_RESYNC string = "FULLRESYNC"
)

func (processor *RespCmdProcessor) ProcessCmd(line string) string {
	parsedResult, err := processor.parser.HandleParse(line)

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "error parsing the line")
	}

	if len(parsedResult) == 0 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	firstCmd := strings.ToUpper(parsedResult[0].Value)
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
	case CMD_INFO:
		return processor.handleInfo(cmds)
	case CMD_REPLCONF:
		return processor.handleReplConf(cmds)
	case CMD_PSYNC:
		return processor.handlePsync(cmds)
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not able to process the cmd")
	}
}
