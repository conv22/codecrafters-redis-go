package cmds

func (processor *RespCmdProcessor) handlePing() (string, error) {
	return processor.parser.HandleEncode(RespEncodingConstants.String, "PONG"), nil
}
