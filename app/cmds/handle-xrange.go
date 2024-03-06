package cmds

import (
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type XRangeHandler struct {
	storage *storage.StorageCollection
}

func newXRangeHandler(storage *storage.StorageCollection) *XRangeHandler {
	return &XRangeHandler{
		storage: storage,
	}
}

func (h *XRangeHandler) minArgs() int {
	return 3
}

func (h *XRangeHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	key, start, end := parsedResult[0].Value, parsedResult[1].Value, parsedResult[2].Value

	value, exist := h.storage.GetCurrentStorage().Get(key)

	if !exist {
		return []string{resp.HandleEncodeSliceList([]resp.SliceEncoding{})}
	}

	stream := value.Value.(*storage.Stream)

	streamRange := stream.GetRange(start, end)

	// todo: move resp conversion to each value type within the storage.
	encoding := strings.Builder{}
	for _, stream := range streamRange {
		encoding.WriteString(resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.BULK_STRING, stream.ID))
		keyValuesEncodings := []resp.SliceEncoding{}

		for key, value := range stream.KeyValues {
			keyValuesEncodings = append(keyValuesEncodings, resp.SliceEncoding{S: key, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING})
			keyValuesEncodings = append(keyValuesEncodings, resp.SliceEncoding{S: value, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING})
		}

		encoding.WriteString(resp.HandleEncodeSliceList(keyValuesEncodings))
	}

	return []string{strconv.Itoa(len(streamRange)) + encoding.String()}

}
