package cmds

import (
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type setKeyOptions struct {
	EX      int   // Set the specified expire time, in seconds (a positive integer).
	PX      int   // Set the specified expire time, in milliseconds (a positive integer).
	EXAT    int64 // Set the specified Unix time at which the key will expire, in seconds (a positive integer).
	PXAT    int64 // Set the specified Unix time at which the key will expire, in milliseconds (a positive integer).
	NX      bool  // Only set the key if it does not already exist.
	XX      bool  // Only set the key if it already exists.
	KEEPTTL bool  // Retain the time to live associated with the key.
}

const (
	EX      = "EX"
	PX      = "PX"
	EXAT    = "EXAT"
	PXAT    = "PXAT"
	NX      = "NX"
	XX      = "XX"
	KEEPTTL = "KEEPTTL"
)

func (processor *RespCmdProcessor) handleSet(parsedResult []resp.ParsedCmd) string {
	if len(parsedResult) < 2 {
		processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}
	key, value := parsedResult[0].Value, parsedResult[1].Value
	var options setKeyOptions
	var expirationTime int64

	if len(parsedResult) >= 3 {
		options = getOptions(parsedResult[2:])
	}

	lockWrite := false

	if options.NX || options.XX {
		_, ok := processor.storage.GetItemFromCurrentStorage(key)

		if !ok && options.XX || ok && options.NX {
			lockWrite = true
		}
	}
	expirationTime = calculateExpirationTime(options)

	if lockWrite {
		return processor.parser.HandleEncode(RespEncodingConstants.NULL_BULK_STRING, "")
	}

	processor.storage.SetItemToCurrentStorage(key, &storage.StorageItem{Value: value, ExpiryMs: expirationTime})
	return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)

}

func getOptions(parsedResult []resp.ParsedCmd) setKeyOptions {
	options := setKeyOptions{}

	for i := 0; i < len(parsedResult); i++ {
		key := strings.ToUpper(parsedResult[i].Value)

		switch key {
		case XX:
			options.XX = true
		case NX:
			options.NX = true
		case KEEPTTL:
			options.KEEPTTL = true
		case EX:
		case PX:
			if i < len(parsedResult)-1 {
				value, err := strconv.Atoi(parsedResult[i+1].Value)
				if err == nil {
					if key == PX {
						options.PX = value
					}
					if key == EX {
						options.EX = value
					}
					i++
				}
			}
		case EXAT:
		case PXAT:
			if i < len(parsedResult)-1 {
				value, err := strconv.ParseInt(parsedResult[i+1].Value, 10, 64)
				if err == nil {
					if key == EXAT {
						options.EXAT = value
					}
					if key == PXAT {
						options.PXAT = value

					}
					i++
				}
			}
		}

	}
	return options
}

func calculateExpirationTime(options setKeyOptions) (timeStamp int64) {
	if options.EXAT > 0 {
		timeStamp = options.EXAT * int64(time.Millisecond)
	}
	if options.PXAT > 0 {
		timeStamp = options.PXAT
	}
	if options.EX > 0 {
		timeStampUnix := time.Now().Add(time.Duration(options.EX) * time.Second).Unix()
		timeStamp = timeStampUnix * int64(time.Millisecond)
	}
	if options.PX > 0 {
		timeStampUnixMilli := time.Now().Add(time.Duration(options.PX) * time.Millisecond).UnixMilli()
		timeStamp = timeStampUnixMilli
	}
	return
}
