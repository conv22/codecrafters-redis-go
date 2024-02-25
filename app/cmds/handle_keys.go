package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type KeysHandler struct {
	storage *storage.StorageCollection
}

func newKeysHandler(storage *storage.StorageCollection) *KeysHandler {
	return &KeysHandler{
		storage: storage,
	}
}

func (h *KeysHandler) minArgs() int {
	return 1
}

func (h *KeysHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	if len(parsedResult) < 1 {
		resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "not enough arguments")
	}

	if parsedResult[0].Value == "*" {

		result := []resp.SliceEncoding{}

		for _, key := range h.storage.GetCurrentStorage().GetKeys() {
			result = append(result, resp.SliceEncoding{S: key, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING})
		}

		return []string{resp.HandleEncodeSliceList(result)}

	}

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, parsedResult[0].Value)}
}
