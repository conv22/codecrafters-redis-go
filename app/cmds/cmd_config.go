package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type CommandHandler interface {
	minArgs() int
	processCmd(parsedCmds []resp.ParsedCmd) []string
}

type PostHandler = func(item string, cmd []resp.ParsedCmd) ProcessCmdResult

type RespCmdProcessor struct {
	handlers     map[string]CommandHandler
	postHandlers map[string]PostHandler
	isMasterConn bool
}

const (
	CMD_PING                 = "PING"
	CMD_ECHO                 = "ECHO"
	CMD_GET                  = "GET"
	CMD_SET                  = "SET"
	CMD_CONFIG               = "CONFIG"
	CMD_KEYS                 = "KEYS"
	CMD_INFO                 = "INFO"
	CMD_REPLCONF             = "REPLCONF"
	CMD_PSYNC                = "PSYNC"
	CMD_RESPONSE_FULL_RESYNC = "FULLRESYNC"
	CMD_RESPONSE_ACK         = "ACK"
	CMD_GETACK               = "GETACK"
	CMD_RESPONSE_OK          = "OK"
	CMD_RESPONSE_PONG        = "PONG"
	CMD_WAIT                 = "WAIT"
	CMD_TYPE                 = "TYPE"
	CMD_XADD                 = "XADD"
)

type ProcessCmdResult struct {
	Answer      string
	BytesInput  []byte
	IsPropagate bool
}
