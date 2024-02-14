package cmds

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

var RespEncodingConstants = parsers.RespEncodingConstants

type RespCmdProcessor struct {
	parser  *parsers.RespParser
	storage *storage.Storage
	config  *config.Config
}

func NewRespCmdProcessor(p *parsers.RespParser, storage *storage.Storage, config *config.Config) *RespCmdProcessor {
	return &RespCmdProcessor{
		parser:  p,
		storage: storage,
		config:  config,
	}
}

const (
	cmdPing   string = "ping"
	cmdEcho   string = "echo"
	cmdGet    string = "get"
	cmdSet    string = "set"
	cmdConfig string = "config"
	cmdKeys   string = "keys"
)

func (processor *RespCmdProcessor) ProcessCmd(line string) string {
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
	case cmdPing:
		return processor.handlePing()
	case cmdEcho:
		return processor.handleEcho(cmds)
	case cmdSet:
		return processor.handleSet(cmds)
	case cmdGet:
		return processor.handleGet(cmds)
	case cmdConfig:
		return processor.handleConfig(cmds)
	case cmdKeys:
		return processor.handleKeys(cmds)
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "not able to process the cmd")
	}
}
