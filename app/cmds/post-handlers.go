package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func noResponsePostHandler(item string, cmds []resp.ParsedCmd) ProcessCmdResult {
	return ProcessCmdResult{
		Answer: "",
	}
}

func propagationPostHandler(item string, cmds []resp.ParsedCmd) ProcessCmdResult {
	return ProcessCmdResult{
		Answer:      item,
		BytesInput:  getBytesInputFromCmds(cmds),
		IsPropagate: true,
	}
}

func defaultPostHandler(item string, cmds []resp.ParsedCmd) ProcessCmdResult {
	return ProcessCmdResult{
		Answer: item,
	}
}
