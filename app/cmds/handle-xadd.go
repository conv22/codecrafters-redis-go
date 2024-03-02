package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

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
	streamEntries := storage.NewStreamEntry(id)
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
