package cmds

import (
	"errors"
	"os"
	"path/filepath"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	reader "github.com/codecrafters-io/redis-starter-go/app/readers"
)

func (processor *RespCmdProcessor) handleKeys(parsedResult []parsers.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}
	reader := reader.CreateReader("rdb")

	dirName, dirFlag := processor.config.DirFlag, processor.config.DbFilenameFlag

	if parsedResult[0].Value == "*" {
		keys, err := reader.HandleRead(filepath.Join(dirName, dirFlag))

		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return processor.parser.HandleEncodeSlice([]parsers.SliceEncoding{})
			}
			return processor.parser.HandleEncode(RespEncodingConstants.Error, err.Error())
		}

		result := []parsers.SliceEncoding{}

		for _, key := range keys {
			result = append(result, parsers.SliceEncoding{Encoding: RespEncodingConstants.BulkString, S: key.Key})
		}

		return processor.parser.HandleEncodeSlice(result)

	}

	return processor.parser.HandleEncode(RespEncodingConstants.String, parsedResult[0].Value)
}
