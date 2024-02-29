package cmds

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

type ReplConfHandler struct {
	cfg *config.Config
}

func newReplConfHandler(cfg *config.Config) *ReplConfHandler {
	return &ReplConfHandler{
		cfg: cfg,
	}
}

func (h *ReplConfHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	firstCmd := parsedResult[0].Value
	switch firstCmd {
	case CMD_GETACK:
		return h.handleGetAck()
	default:
		return h.handleUnknownReplConf()
	}
}

func (h *ReplConfHandler) minArgs() int {
	return 2
}

func (h *ReplConfHandler) handleGetAck() []string {
	return []string{resp.HandleEncodeSliceList([]resp.SliceEncoding{
		{
			S:        CMD_REPLCONF,
			Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING,
		},
		{
			S:        CMD_RESPONSE_ACK,
			Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING,
		},
		{
			S:        strconv.FormatInt(h.cfg.GetOffset(), 10),
			Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING,
		},
	})}
}

func (h *ReplConfHandler) handleUnknownReplConf() []string {
	return []string{resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, CMD_RESPONSE_OK)}
}
