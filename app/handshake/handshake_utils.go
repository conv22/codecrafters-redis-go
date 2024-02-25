package handshake

import (
	"bytes"
	"errors"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func getResponse(conn net.Conn, bufLength int) ([]byte, error) {
	buf := make([]byte, bufLength)

	bytesToRead, err := conn.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf[:bytesToRead], nil
}

func verifyOKResponse(conn net.Conn) error {
	okAnswer := []byte(resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_RESPONSE_OK))
	response, err := getResponse(conn, len(okAnswer))

	if err != nil || !bytes.Equal(okAnswer, response) {
		return errors.New(cmds.CMD_RESPONSE_OK + err.Error())
	}
	return nil
}

func sendOkCommand(conn net.Conn) error {
	return sendCommand(conn, resp.SliceEncoding{S: cmds.CMD_RESPONSE_OK, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING})
}

func sendSliceCommand(conn net.Conn, command []resp.SliceEncoding) error {
	_, err := conn.Write([]byte(resp.HandleEncodeSliceList(command)))
	return err
}

func sendCommand(conn net.Conn, command resp.SliceEncoding) error {
	_, err := conn.Write([]byte(resp.HandleEncode(command.Encoding, command.S)))
	return err
}
