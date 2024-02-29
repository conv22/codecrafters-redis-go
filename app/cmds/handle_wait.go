package cmds

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type WaitHandler struct {
	replicationStore *replication.ReplicationStore
}

func newWaitHandler(replicationStore *replication.ReplicationStore) *WaitHandler {
	return &WaitHandler{
		replicationStore: replicationStore,
	}
}

func (h *WaitHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.INTEGER, strconv.FormatInt(int64(h.replicationStore.NumberOfReplicas()), 10))}
}

func (h *WaitHandler) minArgs() int {
	return 2
}
