package cmds

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

func NewRespCmdProcessor(storage *storage.StorageCollection, config *config.Config, replication *replication.ReplicationStore, conn net.Conn, isMasterConn bool) *RespCmdProcessor {
	processor := &RespCmdProcessor{
		isMasterConn: isMasterConn,
	}

	processor.handlers = make(map[string]CommandHandler)
	processor.postHandlers = make(map[string]PostHandler)

	processor.handlers[CMD_PING] = newPingHandler()
	processor.handlers[CMD_ECHO] = newEchoHandler()
	processor.handlers[CMD_GET] = newGetHandler(storage)
	processor.handlers[CMD_CONFIG] = newConfigHandler(config)
	processor.handlers[CMD_KEYS] = newKeysHandler(storage)
	processor.handlers[CMD_INFO] = newInfoHandler(replication)
	processor.handlers[CMD_TYPE] = newTypeHandler(storage)

	if isMasterConn {
		processor.handlers[CMD_SET] = newSetHandler(storage)
		processor.handlers[CMD_REPLCONF] = newReplConfHandler(replication)
		processor.postHandlers[CMD_REPLCONF] = defaultPostHandler
	} else {
		processor.handlers[CMD_SET] = newSetHandler(storage)
		processor.postHandlers[CMD_SET] = propagationPostHandler
		processor.handlers[CMD_XADD] = newXaddHandler(storage)
		processor.handlers[CMD_XRANGE] = newXRangeHandler(storage)
		processor.postHandlers[CMD_XADD] = propagationPostHandler
		processor.handlers[CMD_PSYNC] = newPsyncHandler(replication, conn)
		processor.handlers[CMD_REPLCONF] = newMasterReplConfHandler(replication, conn)
		processor.handlers[CMD_WAIT] = newWaitHandler(replication)

	}

	return processor
}

func (processor *RespCmdProcessor) ProcessCmd(parsedData []resp.ParsedCmd, conn net.Conn) []ProcessCmdResult {
	result := []ProcessCmdResult{}

	if len(parsedData) == 0 {
		return []ProcessCmdResult{{Answer: resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.ERROR, "not enough arguments")}}
	}

	firstCmd := strings.ToUpper(parsedData[0].Value)
	cmds := parsedData[1:]

	handler, ok := processor.handlers[firstCmd]

	if !ok || len(cmds) < handler.minArgs() {
		return []ProcessCmdResult{}
	}

	postHandler, ok := processor.postHandlers[firstCmd]

	if !ok {
		if processor.isMasterConn {
			postHandler = noResponsePostHandler
		} else {
			postHandler = defaultPostHandler
		}
	}

	for _, item := range handler.processCmd(cmds) {
		result = append(result, postHandler(item, parsedData))
	}

	return result
}
