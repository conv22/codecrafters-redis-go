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
	outputSlices := []resp.SliceEncoding{}

	for _, cmd := range cmds {
		outputSlices = append(outputSlices, resp.SliceEncoding{
			S: cmd.Value, Encoding: cmd.ValueType,
		})
	}

	return ProcessCmdResult{
		Answer:      item,
		BytesInput:  []byte(resp.HandleEncodeSliceList(outputSlices)),
		IsPropagate: true,
	}
}

func defaultPostHandler(item string, cmds []resp.ParsedCmd) ProcessCmdResult {
	return ProcessCmdResult{
		Answer: item,
	}
}
