package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

func (processor *RespCmdProcessor) handleKeys(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	if parsedResult[0].Value == "*" {

		result := []resp.SliceEncoding{}

		for _, key := range processor.storage.GetCurrentStorage().GetKeys() {
			result = append(result, resp.SliceEncoding{S: key, Encoding: RespEncodingConstants.STRING})
		}

		return processor.parser.HandleEncodeSlice(result)

	}

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, parsedResult[0].Value)
}
