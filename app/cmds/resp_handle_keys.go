package cmds

import (
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
)

func (processor *RespCmdProcessor) handleKeys(parsedResult []parsers.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}

	if parsedResult[0].Value == "*" {

		result := []parsers.SliceEncoding{}

		for _, key := range processor.storage.GetKeys() {
			result = append(result, parsers.SliceEncoding{S: key, Encoding: RespEncodingConstants.String})
		}

		return processor.parser.HandleEncodeSlice(result)

	}

	return processor.parser.HandleEncode(RespEncodingConstants.String, parsedResult[0].Value)
}
