package cmds

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type MasterReplConfHandler struct {
	replicationStore *replication.ReplicationStore
	conn             net.Conn
}

func newMasterReplConfHandler(replicationStore *replication.ReplicationStore, conn net.Conn) *MasterReplConfHandler {
	return &MasterReplConfHandler{
		replicationStore: replicationStore,
		conn:             conn,
	}
}

func (h *MasterReplConfHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	replicationAddress, err := replication.GetReplicationAddress(h.conn)
	if err != nil {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "invalid connection address")}
	}

	firstCmd := parsedResult[0].Value
	argument := parsedResult[1].Value

	switch firstCmd {
	case "listening-port":
		return h.handleListeningPort(replicationAddress, argument, h.conn)

	default:
		return h.handleUnknownReplConf(replicationAddress)
	}
}

func (h *MasterReplConfHandler) minArgs() int {
	return 2
}

func (h *MasterReplConfHandler) handleListeningPort(replicationAddress, listeningPort string, conn net.Conn) []string {
	client, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		client = replication.NewReplicaClient(listeningPort)
		h.replicationStore.AppendClient(replicationAddress, client)
	}
	client.AppendConnection(conn)
	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, CMD_RESPONSE_OK)}
}

func (h *MasterReplConfHandler) handleUnknownReplConf(replicationAddress string) []string {
	_, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "invalid handshake")}
	}

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, CMD_RESPONSE_OK)}
}
