package cmds

func (processor *RespCmdProcessor) handlePing() string {
	return processor.parser.HandleEncode(RespEncodingConstants.String, "PONG")
}
