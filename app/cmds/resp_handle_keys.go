package cmds

import (
	"os"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
)

func (processor *RespCmdProcessor) handleKeys(parsedResult []parsers.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}
	dirName, dirFlag := processor.config.DirFlag, processor.config.DbFilenameFlag

	file, err := os.Open(dirName + "/" + dirFlag)
	defer file.Close()

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "no dir provided")
	}
	if parsedResult[0].Value == "*" {

	}

	return processor.parser.HandleEncode(RespEncodingConstants.String, parsedResult[0].Value)
}
