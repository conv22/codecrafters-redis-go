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
	case CMD_RESPONSE_ACK:
		return h.handleAck(replicationAddress, argument, h.conn)
	default:
		return h.handleUnknownReplConf(replicationAddress)
	}
}

func (h *MasterReplConfHandler) minArgs() int {
	return 2
}

func (h *MasterReplConfHandler) handleAck(replicationAddress, offset string, conn net.Conn) []string {
	return []string{}
}

func (h *MasterReplConfHandler) handleListeningPort(replicationAddress, listeningPort string, conn net.Conn) []string {
	client, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		client = replication.NewReplicaClient(listeningPort, conn)
		h.replicationStore.AppendClient(replicationAddress, client)
	}
	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, CMD_RESPONSE_OK)}
}

func (h *MasterReplConfHandler) handleUnknownReplConf(replicationAddress string) []string {
	_, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "invalid handshake")}
	}

	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, CMD_RESPONSE_OK)}
}
