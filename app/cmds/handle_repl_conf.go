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

	firstCmd := parsedResult[0].Value

	switch firstCmd {
	case CMD_GETACK:
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

		// assume that this is the first handshake repl cmd
	case "listening-port":
		client, isClientExist := processor.replication.Replicas[replicationAddress]
		if !isClientExist {
			client = replication.NewReplicaClient(parsedResult[1].Value)
		}
		client.AppendConnection(conn)

		processor.replication.AppendClient(replicationAddress, client)
		return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)

	default:
		_, ok := processor.replication.Replicas[replicationAddress]

		if !ok {
			return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "Invalid handshake")
		}
		return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
	}

}
