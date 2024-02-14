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

	if calculateIsExpired(value.Expiry) {
		processor.storage.GetCurrentStorage().Delete(key)
		return processor.parser.HandleEncode(RespEncodingConstants.NULL_BULK_STRING, "")
	}

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, strValue)

}

func calculateIsExpired(expirationTime *time.Time) bool {
	if expirationTime == nil {
		return false
	}
	return time.Now().After(*expirationTime)
}
