package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type EchoHandler struct{}

func newEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (h *EchoHandler) minArgs() int {
	return 1
}

func (h *EchoHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, parsedResult[0].Value)}
}
