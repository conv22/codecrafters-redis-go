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

type SetHandler struct {
	storage *storage.StorageCollection
}

func newSetHandler(storage *storage.StorageCollection) *SetHandler {
	return &SetHandler{
		storage: storage,
	}
}

func (h *SetHandler) minArgs() int {
	return 2
}

func (h *SetHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	key, value := parsedResult[0].Value, parsedResult[1].Value
	var options setKeyOptions
	var expirationTime int64

	if len(parsedResult) >= 3 {
		options = getOptions(parsedResult[2:])
	}

	lockWrite := false

	if options.NX || options.XX {
		_, ok := h.storage.GetItemFromCurrentStorage(key)

		if !ok && options.XX || ok && options.NX {
			lockWrite = true
		}
	}
	expirationTime = calculateExpirationTime(options)

	if lockWrite {
		return []string{resp.HandleEncode(respEncodingConstants.NULL_BULK_STRING, "")}
	}

	h.storage.SetItemToCurrentStorage(key, &storage.StorageItem{Value: value, ExpiryMs: expirationTime})

	return []string{resp.HandleEncode(respEncodingConstants.STRING, CMD_RESPONSE_OK)}

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
		case EX, PX:
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
		case EXAT, PXAT:
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
