package cmds

import (
	"errors"
	"strings"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

var RespEncodingConstants = parsers.RespEncodingConstants

type RespCmdProcessor struct {
	parser  parsers.Parser
	storage storage.Storage
}

func (processor RespCmdProcessor) ProcessCmd(line string) (string, error) {
	parsedResult, err := processor.parser.HandleParse(line)

	if err != nil {
		return "", err
	}

	if len(parsedResult) == 0 {
		return "", errors.New("no arguments were parsed")
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
	default:
		return "", errors.New("not able to process cmd")
	}
}
