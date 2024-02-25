package cmds

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func NewRespCmdProcessor(storage *storage.StorageCollection, config *config.Config, replication *replication.ReplicationStore, conn net.Conn) *RespCmdProcessor {
	processor := &RespCmdProcessor{}

	processor.handlers = make(map[string]CommandHandler)
	processor.postHandlers = make(map[string]PostHandler)

	processor.handlers[CMD_PING] = newPingHandler()
	processor.handlers[CMD_ECHO] = newEchoHandler()
	processor.handlers[CMD_GET] = newGetHandler(storage)
	processor.handlers[CMD_CONFIG] = newConfigHandler(config)
	processor.handlers[CMD_KEYS] = newKeysHandler(storage)
	processor.handlers[CMD_INFO] = newInfoHandler(replication)

	if replication.IsReplica() {
		// ignore set commands from other clients.
		if replication.IsCmdFromMaster(conn) {
			processor.handlers[CMD_SET] = newSetHandler(storage)
			processor.postHandlers[CMD_SET] = noResponsePostHandler
		}
	} else {
		processor.handlers[CMD_SET] = newSetHandler(storage)
		processor.postHandlers[CMD_SET] = propagationPostHandler
		processor.handlers[CMD_PSYNC] = newPsyncHandler(replication, conn)
		processor.handlers[CMD_REPLCONF] = newReplConfHandler(replication, conn)
	}

	return processor
}

func getBytesInputFromCmds(cmds []resp.ParsedCmd) []byte {
	outputSlices := []resp.SliceEncoding{}

	for _, cmd := range cmds {
		outputSlices = append(outputSlices, resp.SliceEncoding{
			S: cmd.Value, Encoding: cmd.ValueType,
		})
	}

	return []byte(resp.HandleEncodeSliceList(outputSlices))

}

func (processor *RespCmdProcessor) ProcessCmd(data []byte, conn net.Conn) []ProcessCmdResult {
	result := []ProcessCmdResult{}
	parsedLines, err := resp.HandleParse(string(data))

	if err != nil {
		return []ProcessCmdResult{{Answer: resp.HandleEncode(respEncodingConstants.ERROR, "error parsing the line")}}
	}

	if len(parsedLines) == 0 {
		return []ProcessCmdResult{{Answer: resp.HandleEncode(respEncodingConstants.ERROR, "not enough arguments")}}
	}

	for _, parsedLine := range parsedLines {

		firstCmd := strings.ToUpper(parsedLine[0].Value)
		cmds := parsedLine[1:]

		handler, ok := processor.handlers[firstCmd]

		if !ok || len(cmds) < handler.minArgs() {
			continue
		}

		postHandler, ok := processor.postHandlers[firstCmd]

		if !ok {
			postHandler = defaultPostHandler
		}

		for _, item := range handler.processCmd(cmds) {
			result = append(result, postHandler(item, parsedLine))
		}
	}

	return result
}
