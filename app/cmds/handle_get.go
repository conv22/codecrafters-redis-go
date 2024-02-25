package cmds

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type GetHandler struct {
	storage *storage.StorageCollection
}

func newGetHandler(storage *storage.StorageCollection) *GetHandler {
	return &GetHandler{
		storage: storage,
	}
}

func (h *GetHandler) minArgs() int {
	return 1
}

func (h *GetHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	key := parsedResult[0].Value
	value, ok := h.storage.GetCurrentStorage().Get(key)
	if !ok {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.NULL_BULK_STRING, "")}
	}
	strValue, ok := value.Value.(string)
	if !ok {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.NULL_BULK_STRING, "")}
	}

	if calculateIsExpired(value.ExpiryMs) {
		h.storage.GetCurrentStorage().Delete(key)
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.NULL_BULK_STRING, "")}
	}

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, strValue)}

}

func calculateIsExpired(expirationTime int64) bool {
	if expirationTime == 0 {
		return false
	}
	currentTime := time.Now().UnixMilli()
	return currentTime > expirationTime
}
