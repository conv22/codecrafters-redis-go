package cmds

import "github.com/codecrafters-io/redis-starter-go/app/resp"

type PingHandler struct{}

func newPingHandler() *PingHandler {
	return &PingHandler{}
}

func (p *PingHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	return []string{resp.HandleEncode(respEncodingConstants.STRING, CMD_PONG)}
}

func (h *PingHandler) minArgs() int {
	return 0
}
