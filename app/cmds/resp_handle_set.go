package cmds

import (
	"errors"
	"strconv"
	"strings"
	"time"

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
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

func (processor *RespCmdProcessor) handleSet(parsedResult []parsers.ParsedCmd) (string, error) {
	if len(parsedResult) < 2 {
		return "", errors.New("not enough arguments")
	}
	key, value := parsedResult[0].Value, parsedResult[1].Value
	var options *setKeyOptions = nil

	if len(parsedResult) >= 3 {
		options = getOptions(parsedResult[2:])
	}

	lockWrite := false

	if options != nil && (options.NX || options.XX) {
		currentValue := processor.storage.Get(storage.StorageKey{Key: key})

		if currentValue == nil && options.XX || currentValue != nil && options.NX {
			lockWrite = true
		}
	}

	if lockWrite {
		return processor.parser.HandleEncode(RespEncodingConstants.NullBulkString, ""), nil
	}

	processor.storage.Set(storage.StorageKey{Key: key}, storage.StorageValue{Value: value, ExpirationTime: calculateExpirationTime(options)})
	return processor.parser.HandleEncode(RespEncodingConstants.String, "OK"), nil

}

func getOptions(parsedResult []parsers.ParsedCmd) *setKeyOptions {
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
	return &options
}

func calculateExpirationTime(options *setKeyOptions) *int64 {
	if options == nil {
		return nil
	}
	if options.EXAT > 0 {
		return &options.EXAT
	}
	if options.PXAT > 0 {
		pxatSeconds := options.PXAT / 1000
		return &pxatSeconds
	}
	if options.EX > 0 {
		exTime := time.Now().Add(time.Duration(options.EX) * time.Second).UnixMilli()
		return &exTime
	}
	if options.PX > 0 {
		pxTime := time.Now().Add(time.Duration(options.PX) * time.Millisecond).UnixMilli()
		return &pxTime
	}
	return nil
}
