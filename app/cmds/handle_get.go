package cmds

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (processor *RespCmdProcessor) handleGet(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}
	key := parsedResult[0].Value
	value, ok := processor.storage.Get(key)
	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.NullBulkString, "")
	}
	if calculateIsExpired(value.Expiry) {
		processor.storage.Delete(key)
		return processor.parser.HandleEncode(RespEncodingConstants.NullBulkString, "")
	}
	return processor.parser.HandleEncode(RespEncodingConstants.String, value.Value.(string))

}

func calculateIsExpired(expirationTime *time.Time) bool {
	if expirationTime == nil {
		return false
	}
	return time.Now().After(*expirationTime)
}
