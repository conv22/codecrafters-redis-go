package cmds

import (
	"flag"
	"strings"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
)

func (processor RespCmdProcessor) handleConfig(parsedResult []parsers.ParsedCmd) string {
	if len(parsedResult) < 2 {
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}
	cmd := strings.ToLower(parsedResult[0].Value)

	switch cmd {
	case "get":
		{
			flagType := parsedResult[1].Value
			value := ""
			if flagType == "dir" {
				dirFlag := processor.config.DirFlag
				flag.Parse()
				value = dirFlag

			}

			if flagType == "dbfilename" {
				dbFileNameFlag := processor.config.DbFilenameFlag
				flag.Parse()
				value = dbFileNameFlag
			}

			encodings := []parsers.SliceEncoding{
				{S: flagType, Encoding: RespEncodingConstants.BulkString},
				{S: value, Encoding: RespEncodingConstants.BulkString},
			}

			return processor.parser.HandleEncodeSlice(encodings)
		}
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.Error, "unsupported cmd")
	}

}
