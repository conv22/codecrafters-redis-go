package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type PingHandler struct{}

func newPingHandler() *PingHandler {
	return &PingHandler{}
}

func (p *PingHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, CMD_RESPONSE_PONG)}
}

func (h *PingHandler) minArgs() int {
	return 0
}
