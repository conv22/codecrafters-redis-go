package cmds

import (
	"errors"
	"time"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

func (processor *RespCmdProcessor) handleGet(parsedResult []parsers.ParsedCmd) (string, error) {
	if len(parsedResult) < 1 {
		return "", errors.New("not enough arguments")
	}
	key := parsedResult[0].Value
	value := processor.storage.Get(storage.StorageKey{Key: key})
	if value == nil {
		return "", nil
	}
	if calculateIsExpired(value.ExpirationTime) {
		processor.storage.Delete(storage.StorageKey{Key: key})
		return processor.parser.HandleEncode(RespEncodingConstants.NullBulkString, ""), nil
	}
	return processor.parser.HandleEncode(RespEncodingConstants.String, value.Value), nil

}

func calculateIsExpired(unix *int64) bool {
	if unix == nil {
		return false
	}
	expirationTime := time.UnixMilli(*unix)

	return time.Now().After(expirationTime)

}
