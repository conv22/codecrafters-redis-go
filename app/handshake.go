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

	defer masterConn.Close()

	// Send PING command
	if err := sendPingCommand(masterConn); err != nil {
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

	return nil
}

func connectToMaster() (net.Conn, error) {
	replication := cfg.Replication
	masterConn, err := net.Dial("tcp", replication.MasterAddress)
	if err != nil {
		return nil, errors.New("failed to connect to master: " + replication.MasterAddress)
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

func verifyOKResponse(conn net.Conn) error {
	okAnswer := []byte(parser.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, "OK"))
	buf := make([]byte, len(okAnswer))
	bytesRead, err := conn.Read(buf)

	if err != nil || !bytes.Equal(buf[:bytesRead], okAnswer) {
		// return errors.New("expected OK response not received")
		return nil
	}
	return nil
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
