package main

import (
	"bufio"
	"bytes"
	"errors"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

func handleHandshake() error {
	// Establish connection to the master
	masterConn, err := connectToMaster()
	if err != nil {
		return err
	}

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

func connectToMaster() (net.Conn, error) {
	masterConn, err := net.Dial("tcp", replicationInfo.MasterAddress)
	if err != nil {
		return nil, errors.New("failed to connect to master: " + replicationInfo.MasterAddress)
	}
	return masterConn, nil
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
		return errors.New("failed to send PSYNC command: " + err.Error())
	}
	return nil
}

func sendListeningPortConfig(conn net.Conn) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	listeningPortConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "listening-port", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: cfg.Port, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	if err := sendCommand(writer, listeningPortConfig); err != nil {
		return errors.New("failed to send listening port configuration: " + err.Error())
	}
	return nil
}

func verifyPingResponse(conn net.Conn) error {
	pingAnswer := []byte(parser.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_PONG))
	response, err := getResponse(conn, len(pingAnswer))

	if err != nil || !bytes.Equal(pingAnswer, response) {
		return errors.New(cmds.CMD_PONG + err.Error())
	}
	return nil
}

func verifyOKResponse(conn net.Conn) error {
	okAnswer := []byte(parser.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_OK))
	response, err := getResponse(conn, len(okAnswer))

	if err != nil || !bytes.Equal(okAnswer, response) {
		return errors.New(cmds.CMD_OK + err.Error())
	}
	return nil
}

func verifyPsyncResponse(conn net.Conn) error {
	_, err := getResponse(conn, 1024)

	if err != nil {
		return errors.New(cmds.CMD_PSYNC + err.Error())
	}
	return nil
}

func getResponse(conn net.Conn, bufLength int) ([]byte, error) {
	buf := make([]byte, bufLength)

	bytesToRead, err := conn.Read(buf)

	if err != nil {
		return nil, errors.New("expected OK response not received")
	}

	return buf[0:bytesToRead], nil
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
	_, err := writer.Write([]byte(parser.HandleEncodeSliceList(command)))
	return err
}
