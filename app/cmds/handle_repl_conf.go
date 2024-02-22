package cmds

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (processor *RespCmdProcessor) handleReplConf(parsedResult []resp.ParsedCmd, conn net.Conn) string {
	// conn, parse and save port
	if len(parsedResult) < 2 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	replicationAddress, err := replication.GetReplicationAddress(conn)

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "Invalid connection address")
	}

	// assume that this is the first handshake repl cmd
	if parsedResult[0].Value == "listening-port" {
		client, isClientExist := processor.replication.Replicas[replicationAddress]
		if !isClientExist {
			client = replication.NewReplicaClient(parsedResult[1].Value)
		}
		client.AppendConnection(conn)

		processor.replication.AppendClient(replicationAddress, client)
		return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
	}

	replica, ok := processor.replication.Replicas[replicationAddress]

	if !ok {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "Invalid handshake")
	}

	// assume that this is the second handshake repl cmd, ignore for now as the values are hardcoded

	// handshake complete
	replica.IsActive = true

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
}
