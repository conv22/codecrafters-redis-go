package cmds

import (
	"errors"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

var errTooSmallStreamId = errors.New("ERR The ID specified in XADD is equal or smaller than the target stream top item")
var errStreamIdZero = errors.New("ERR The ID specified in XADD must be greater than 0-0")
var errStreamIdFormat = errors.New("ERR the ID specified in XADD must be in  milliseconds-sqnumber format")

type XaddHandler struct {
	storage *storage.StorageCollection
}

func newXaddHandler(storage *storage.StorageCollection) *XaddHandler {
	return &XaddHandler{
		storage: storage,
	}
}

func (h *XaddHandler) minArgs() int {
	return 3
}

func (h *XaddHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	key, id := parsedResult[0].Value, h.getNewStreamEntryId(parsedResult[1].Value)
	entries := parsedResult[2:]
	stream, ok := h.storage.GetCurrentStorage().Get(key)

	var currentStream *storage.Stream

	if !ok {
		currentStream = storage.NewStream()
	} else {
		currentStream = stream.Value.(*storage.Stream)
	}

	// new or find, todo
	msTime, sqNumber, err := getParsedStreamId(id, currentStream)
	if err != nil {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, err.Error())}
	}
	streamEntries := storage.NewStreamEntry(msTime, sqNumber, currentStream)

	for i := 0; i < len(entries); i += 2 {
		if i+1 >= len(entries) {
			break
		}
		key, value := entries[i].Value, entries[i+1].Value
		streamEntries.AddEntry(key, value)
	}

	currentStream.AddEntry(streamEntries)
	h.storage.SetItemToCurrentStorage(key, &storage.StorageItem{
		Type:  storage.STREAM,
		Value: currentStream,
	})

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.BULK_STRING, id)}
}
func (h *XaddHandler) getNewStreamEntryId(id string) string {
	if id == "*" {
		// autogenerate later
		return id
	}
	return id
}

func getParsedStreamId(id string, stream *storage.Stream) (msTime, sqNumber int64, err error) {
	msTime, sqNumber, err = getMsAndSqFromId(id, stream)

	if err != nil {
		return 0, 0, errStreamIdZero
	}

	if msTime == 0 && sqNumber == 0 {
		return 0, 0, errStreamIdZero
	}
	if len(stream.Entries) == 0 {
		return msTime, sqNumber, nil
	}

	lastEntry := stream.Entries[len(stream.Entries)-1]

	if lastEntry.MsTime > msTime || lastEntry.MsTime == msTime && lastEntry.SqNumber >= sqNumber {
		return 0, 0, errTooSmallStreamId
	}

	return msTime, sqNumber, nil
}

func getMsAndSqFromId(id string, stream *storage.Stream) (msTime, sqNumber int64, err error) {
	split := strings.Split(id, "-")
	if len(split) != 2 {
		return 0, 0, errStreamIdFormat
	}

	msTime, err = strconv.ParseInt(split[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	sqNumber, err = getSqNumber(split[1], msTime, stream)
	if err != nil {
		return 0, 0, err
	}

	return msTime, sqNumber, nil

}

func getSqNumber(sq string, msTime int64, stream *storage.Stream) (sqNumber int64, err error) {
	if sq == "*" {
		if len(stream.Entries) == 0 {
			if msTime == 0 {
				return 1, nil
			} else {
				return 0, nil
			}
		}
		lastEntry := stream.Entries[len(stream.Entries)-1]
		return lastEntry.SqNumber + 1, nil

	}
	// explicit
	sqNumber, err = strconv.ParseInt(sq, 10, 64)
	if err != nil {
		return 0, err
	}

	return sqNumber, nil
}
