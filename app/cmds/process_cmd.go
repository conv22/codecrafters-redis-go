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

func (processor *RespCmdProcessor) ProcessCmd(line string) (str string, strSlice []string, isSlice bool) {
	parsedResult, err := processor.parser.HandleParse(line)

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "error parsing the line"), nil, false
	}

	if len(parsedResult) == 0 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments"), nil, false
	}

	firstCmd := strings.ToUpper(parsedResult[0].Value)
	cmds := parsedResult[1:]

	switch firstCmd {
	case CMD_PING:
		return processor.handlePing(), nil, false
	case CMD_ECHO:
		return processor.handleEcho(cmds), nil, false
	case CMD_SET:
		return processor.handleSet(cmds), nil, false
	case CMD_GET:
		return processor.handleGet(cmds), nil, false
	case CMD_CONFIG:
		return processor.handleConfig(cmds), nil, false
	case CMD_KEYS:
		return processor.handleKeys(cmds), nil, false
	case CMD_INFO:
		return processor.handleInfo(cmds), nil, false
	case CMD_REPLCONF:
		return processor.handleReplConf(cmds), nil, false
	case CMD_PSYNC:
		return "", processor.handlePsync(cmds), true
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not able to process the cmd"), nil, false
	}
}
