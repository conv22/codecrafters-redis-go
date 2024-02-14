package cmds

import (
	"time"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

func (processor *RespCmdProcessor) handleGet(parsedResult []parsers.ParsedCmd) string {
	if len(parsedResult) < 1 {
		processor.parser.HandleEncode(RespEncodingConstants.Error, "not enough arguments")
	}
	key := parsedResult[0].Value
	value, ok := processor.storage.Get(storage.StorageKey{Key: key})
	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.NullBulkString, "")
	}
	if calculateIsExpired(value.ExpirationTime) {
		processor.storage.Delete(storage.StorageKey{Key: key})
		return processor.parser.HandleEncode(RespEncodingConstants.NullBulkString, "")
	}
	return processor.parser.HandleEncode(RespEncodingConstants.String, value.Value.(string))

}

func calculateIsExpired(unix *int64) bool {
	if unix == nil {
		return false
	}
	expirationTime := time.UnixMilli(*unix)

	return time.Now().After(expirationTime)

}
