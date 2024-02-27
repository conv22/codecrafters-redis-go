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
}

const (
	CMD_PING          = "PING"
	CMD_ECHO          = "ECHO"
	CMD_GET           = "GET"
	CMD_SET           = "SET"
	CMD_CONFIG        = "CONFIG"
	CMD_KEYS          = "KEYS"
	CMD_INFO          = "INFO"
	CMD_RESPONSE_OK   = "OK"
	CMD_RESPONSE_PONG = "PONG"
)

type ProcessCmdResult struct {
	Answer      string
	BytesInput  []byte
	IsPropagate bool
}