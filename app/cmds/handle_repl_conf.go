package cmds

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ReplConfHandler struct {
	replicationStore *replication.ReplicationStore
	conn             net.Conn
}

func newReplConfHandler(replicationStore *replication.ReplicationStore, conn net.Conn) *ReplConfHandler {
	return &ReplConfHandler{
		replicationStore: replicationStore,
		conn:             conn,
	}
}

func (h *ReplConfHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	replicationAddress, err := replication.GetReplicationAddress(h.conn)
	if err != nil {
		return []string{resp.HandleEncode(respEncodingConstants.ERROR, "invalid connection address")}
	}

	firstCmd := parsedResult[0].Value
	argument := parsedResult[1].Value

	switch firstCmd {
	case CMD_RESPONSE_ACK:
		return h.handleGetAck(replicationAddress)

	case "listening-port":
		return h.handleListeningPort(replicationAddress, argument, h.conn)

	default:
		return h.handleUnknownReplConf(replicationAddress)
	}
}

func (h *ReplConfHandler) minArgs() int {
	return 2
}

func (h *ReplConfHandler) handleGetAck(replicationAddress string) []string {
	_, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return []string{resp.HandleEncode(respEncodingConstants.ERROR, "invalid handshake")}
	}

	return []string{resp.HandleEncodeSliceList([]resp.SliceEncoding{
		{
			S:        CMD_REPLCONF,
			Encoding: respEncodingConstants.BULK_STRING,
		},
		{
			S:        CMD_RESPONSE_ACK,
			Encoding: respEncodingConstants.BULK_STRING,
		},
		{
			S:        "0",
			Encoding: respEncodingConstants.BULK_STRING,
		},
	})}
}

func (h *ReplConfHandler) handleListeningPort(replicationAddress, listeningPort string, conn net.Conn) []string {
	client, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		client = replication.NewReplicaClient(listeningPort)
		h.replicationStore.AppendClient(replicationAddress, client)
	}

	client.AppendConnection(conn)
	return []string{resp.HandleEncode(respEncodingConstants.STRING, CMD_RESPONSE_OK)}
}

func (h *ReplConfHandler) handleUnknownReplConf(replicationAddress string) []string {
	_, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return []string{resp.HandleEncode(respEncodingConstants.ERROR, "invalid handshake")}
	}

	return []string{resp.HandleEncode(respEncodingConstants.STRING, CMD_RESPONSE_OK)}
}
