package cmds

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (processor *RespCmdProcessor) handleReplConf(parsedResult []resp.ParsedCmd, conn net.Conn) string {
	if len(parsedResult) < 2 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	replicationAddress, err := replication.GetReplicationAddress(conn)
	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "invalid connection address")
	}

	firstCmd := parsedResult[0].Value
	argument := parsedResult[1].Value

	switch firstCmd {
	case CMD_GETACK:
		return processor.handleGetAck(replicationAddress)

	case "listening-port":
		return processor.handleListeningPort(replicationAddress, argument, conn)

	default:
		return processor.handleUnknownReplConf(replicationAddress)
	}
}

func (processor *RespCmdProcessor) handleGetAck(replicationAddress string) string {
	_, ok := processor.replication.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "invalid handshake")
	}

	return processor.parser.HandleEncodeSliceList([]resp.SliceEncoding{
		{
			S:        CMD_REPLCONF,
			Encoding: RespEncodingConstants.BULK_STRING,
		},
		{
			S:        CMD_ACK,
			Encoding: RespEncodingConstants.BULK_STRING,
		},
		{
			S:        "0",
			Encoding: RespEncodingConstants.BULK_STRING,
		},
	})
}

func (processor *RespCmdProcessor) handleListeningPort(replicationAddress, listeningPort string, conn net.Conn) string {
	client, ok := processor.replication.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		client = replication.NewReplicaClient(listeningPort)
		processor.replication.AppendClient(replicationAddress, client)
	}

	client.AppendConnection(conn)
	return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
}

func (processor *RespCmdProcessor) handleUnknownReplConf(replicationAddress string) string {
	_, ok := processor.replication.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "invalid handshake")
	}

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
}
