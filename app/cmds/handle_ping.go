package cmds

func (processor *RespCmdProcessor) handlePing() string {
	return processor.parser.HandleEncode(RespEncodingConstants.STRING, "PONG")
}
