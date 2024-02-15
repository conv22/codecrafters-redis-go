package cmds

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (processor *RespCmdProcessor) handleGet(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}
	key := parsedResult[0].Value
	value, ok := processor.storage.GetCurrentStorage().Get(key)
	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.NULL_BULK_STRING, "")
	}
	strValue, ok := value.Value.(string)
	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.NULL_BULK_STRING, "")
	}

	if calculateIsExpired(value.ExpiryMs) {
		processor.storage.GetCurrentStorage().Delete(key)
		return processor.parser.HandleEncode(RespEncodingConstants.NULL_BULK_STRING, "")
	}

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, strValue)

}

func calculateIsExpired(expirationTime int64) bool {
	if expirationTime == 0 {
		return false
	}
	currentTime := time.Now().UnixMilli()
	return currentTime > expirationTime
}
