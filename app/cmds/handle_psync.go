package cmds

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type PsyncHandler struct {
	replicationStore *replication.ReplicationStore
	conn             net.Conn
}

func newPsyncHandler(replicationStore *replication.ReplicationStore, conn net.Conn) *PsyncHandler {
	return &PsyncHandler{
		replicationStore: replicationStore,
		conn:             conn,
	}
}

func (h *PsyncHandler) minArgs() int {
	return 2
}

const EMPTY_DB_HEX string = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

func (h *PsyncHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	replicationAddress, err := replication.GetReplicationAddress(h.conn)

	if err != nil {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "Invalid connection address")}
	}

	replica, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)

	if !ok {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "Invalid connection address")}
	}

	offset, replicationId := parsedResult[0], parsedResult[1]
	replica.SetOffsetAndReplicationId(offset.Value, replicationId.Value)

	decoded, err := hex.DecodeString(EMPTY_DB_HEX)

	if err != nil {
		return nil
	}

	ackCmd := resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, fmt.Sprintf("%s %s %s", CMD_RESPONSE_FULL_RESYNC, h.replicationStore.MasterReplId, h.replicationStore.Offset))

	encodingCmd := resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.BULK_STRING, string(decoded))

	// exception
	encodingCmd = strings.TrimSuffix(encodingCmd, resp.RESP_ENCODING_CONSTANTS.SEPARATOR)

	return []string{ackCmd, encodingCmd}
}
