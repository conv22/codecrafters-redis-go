package cmds

import (
	"strconv"
	"time"

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

func (h *WaitHandler) minArgs() int {
	return 2
}

func (h *WaitHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	expectedReplicasStr := parsedResult[0].Value
	_, err := strconv.ParseInt(expectedReplicasStr, 10, 32)

	if err != nil {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "Invalid timeout passed")}
	}

	timeStr := parsedResult[1].Value
	timeInt, err := strconv.ParseInt(timeStr, 10, 32)

	if err != nil {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "Invalid timeout passed")}
	}

	if !h.replicationStore.HasReplicas() {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.INTEGER, "0")}
	}

	acknowledgedCount := h.replicationStore.GetNumOfAckReplicas()

	if acknowledgedCount == h.replicationStore.NumberOfReplicas() {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.INTEGER, strconv.FormatInt(int64(h.replicationStore.NumberOfReplicas()), 10))}
	}

	go h.replicationStore.GetAckFromReplicas()

	acknowledgedCh := make(chan struct{})
	timer := time.NewTimer(time.Duration(timeInt) * time.Millisecond)

	go func() {
	outer:
		for {
			select {
			case <-h.replicationStore.AckCompleteChan:
				acknowledgedCount++
			case <-timer.C:
				close(acknowledgedCh)
				break outer
			}
		}
	}()

	<-acknowledgedCh

	h.replicationStore.GetNumOfAckReplicas()

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.INTEGER, strconv.FormatInt(int64(acknowledgedCount), 10))}

}
