package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func (processor *RespCmdProcessor) handleReplConf(parsedResult []resp.ParsedCmd) string {
	// conn, parse and save port
	if len(parsedResult) < 2 {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")
	}

	replicationAddress, err := replication.GetReplicationAddress(processor.connection)

	if err != nil {
		return processor.parser.HandleEncode(RespEncodingConstants.ERROR, "Invalid connection address")
	}

	replica, ok := processor.replication.Replicas[replicationAddress]

	// assume that this is the first handshake repl cmd
	if !ok {
		port := parsedResult[1].Value
		processor.replication.Replicas[replicationAddress] = replication.NewReplicaClient(processor.connection, port)
		return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
	}
	// assume that this is the second handshake repl cmd, ignore for now as the values are hardcoded

	// handshake complete
	replica.IsActive = true

	return processor.parser.HandleEncode(RespEncodingConstants.STRING, CMD_OK)
}
