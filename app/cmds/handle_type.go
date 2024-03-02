package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type TypeHandler struct {
	storage *storage.StorageCollection
}

func newTypeHandler(storage *storage.StorageCollection) *TypeHandler {
	return &TypeHandler{
		storage: storage,
	}
}

func (h *TypeHandler) minArgs() int {
	return 1
}

func (h *TypeHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	key := parsedResult[0].Value
	value, ok := h.storage.GetCurrentStorage().Get(key)

	if !ok {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, storage.NoneType)}
	}

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, value.Type)}

}
