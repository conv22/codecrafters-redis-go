package cmds

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type InfoHandler struct {
	replicationStore *replication.ReplicationStore
}

func newInfoHandler(replicationStore *replication.ReplicationStore) *InfoHandler {
	return &InfoHandler{
		replicationStore: replicationStore,
	}
}

func (h *InfoHandler) minArgs() int {
	return 1
}

const (
	INFO_CMD_REPLICATION = "replication"
)

func (h *InfoHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	switch parsedResult[0].Value {
	case INFO_CMD_REPLICATION:
		replication := h.replicationStore
		data := []resp.SliceEncoding{
			{S: "role:" + replication.Role, Encoding: resp.RESP_ENCODING_CONSTANTS.SEPARATOR},
			{S: "master_replid:" + replication.MasterReplId, Encoding: resp.RESP_ENCODING_CONSTANTS.SEPARATOR},
			{S: "master_repl_offset:" + strconv.FormatInt(h.replicationStore.GetOffset(), 10), Encoding: resp.RESP_ENCODING_CONSTANTS.SEPARATOR},
		}

		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.BULK_STRING, resp.HandleEncodeSlices(data))}
	default:
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "invalid argument")}
	}
}
