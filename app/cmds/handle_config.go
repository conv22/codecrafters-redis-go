package cmds

import (
	"flag"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	configDir      = "dir"
	configFileName = "dbfilename"
)

func (processor *RespCmdProcessor) handleConfig(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 2 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}
	cmd := strings.ToLower(parsedResult[0].Value)

	switch cmd {
	case "get":
		{
			flagType := parsedResult[1].Value
			value := ""
			if flagType == configDir {
				dirFlag := processor.config.DirFlag
				flag.Parse()
				value = dirFlag

			}

			if flagType == configFileName {
				dbFileNameFlag := processor.config.DbFilenameFlag
				flag.Parse()
				value = dbFileNameFlag
			}

			encodings := []resp.SliceEncoding{
				{S: flagType, Encoding: RespEncodingConstants.BULK_STRING},
				{S: value, Encoding: RespEncodingConstants.BULK_STRING},
			}

			return processor.parser.HandleEncodeSlice(encodings)
		}
	default:
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "unsupported cmd")
	}

}
