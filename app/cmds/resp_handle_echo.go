package cmds

import (
	"errors"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
)

func (processor *RespCmdProcessor) handleEcho(parsedResult []parsers.ParsedCmd) (string, error) {
	if len(parsedResult) < 1 {
		return "", errors.New("not enough arguments")

	}
	return processor.parser.HandleEncode(RespEncodingConstants.String, parsedResult[0].Value), nil
}
