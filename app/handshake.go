package main

import (
	"bufio"
	"bytes"
	"errors"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handleHandshake(masterConn net.Conn) error {
	// Send PING command
	if err := sendPingCommand(masterConn); err != nil {
		return err
	}

	// Verify PING response from master
	if err := verifyPingResponse(masterConn); err != nil {
		return err
	}

	// Send listening port configuration
	if err := sendListeningPortConfig(masterConn); err != nil {
		return err
	}

	// Verify OK response for listening port configuration
	if err := verifyOKResponse(masterConn); err != nil {
		return err
	}

	// Send capability configuration
	if err := sendCapabilityConfig(masterConn); err != nil {
		return err
	}

	// Verify OK response for capability configuration
	if err := verifyOKResponse(masterConn); err != nil {
		return err
	}

	// Send psync

	if err := sendPsyncCommand(masterConn); err != nil {
		return err
	}

	// Verify PSYNC response
	if err := verifyPsyncResponse(masterConn); err != nil {
		return err
	}

	return nil
}

func sendPingCommand(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	pingCommand := []resp.SliceEncoding{{S: cmds.CMD_PING, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING}}

	if err := sendCommand(writer, pingCommand); err != nil {
		return errors.New("failed to send PING command: " + err.Error())
	}
	return nil
}

func sendPsyncCommand(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	psyncCommand := []resp.SliceEncoding{
		{S: cmds.CMD_PSYNC, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "?", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "-1", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	if err := sendCommand(writer, psyncCommand); err != nil {
		return err
	}

	return nil
}

func sendListeningPortConfig(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	listeningPortConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "listening-port", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: serverContext.cfg.Port, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	if err := sendCommand(writer, listeningPortConfig); err != nil {
		return errors.New("failed to send listening port configuration: " + err.Error())
	}
	return nil
}

func verifyPingResponse(conn net.Conn) error {
	pingAnswer := []byte(serverContext.parser.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_PONG))
	response, err := getResponse(conn, len(pingAnswer))

	if err != nil || !bytes.Equal(pingAnswer, response) {
		return errors.New(cmds.CMD_PONG + err.Error())
	}
	return nil
}

func verifyOKResponse(conn net.Conn) error {
	okAnswer := []byte(serverContext.parser.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_OK))
	response, err := getResponse(conn, len(okAnswer))

	if err != nil || !bytes.Equal(okAnswer, response) {
		return errors.New(cmds.CMD_OK + err.Error())
	}
	return nil
}

func verifyPsyncResponse(conn net.Conn) error {
	// full resync
	_, err := getResponse(conn, 1024)

	if err != nil {
		return errors.New(cmds.CMD_PSYNC + err.Error())
	}

	// // rdb file
	// rdbFileBytes, err := getResponse(conn, 1024)

	// if err != nil {
	// 	collection, err := serverContext.rdbReader.HandleReadFromBytes(rdbFileBytes)
	// 	if err != nil {
	// 		serverContext.inMemoryStorage = collection
	// 	}
	// }

	return err
}

func getResponse(conn net.Conn, bufLength int) ([]byte, error) {
	buf := make([]byte, bufLength)

	bytesToRead, err := conn.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf[:bytesToRead], nil
}

func sendCapabilityConfig(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	capabilityConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "capa", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "psync2", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	if err := sendCommand(writer, capabilityConfig); err != nil {
		return errors.New("failed to send capability configuration: " + err.Error())
	}
	return nil
}

func sendCommand(writer *bufio.Writer, command []resp.SliceEncoding) error {
	_, err := writer.Write([]byte(serverContext.parser.HandleEncodeSliceList(command)))
	return err
}
