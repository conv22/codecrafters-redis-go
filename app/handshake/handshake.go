package handshake

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type handshake struct {
	cfg             *config.Config
	rdbReader       *rdb.Rdb
	inMemoryStorage *storage.StorageCollection
}

func New(cfg *config.Config, inMemoryStorage *storage.StorageCollection) *handshake {
	return &handshake{
		cfg:             cfg,
		rdbReader:       rdb.NewRdb(),
		inMemoryStorage: inMemoryStorage,
	}
}

func (h *handshake) HandleHandshake(masterConn net.Conn) error {
	var steps []func() error

	steps = append(steps, func() error {
		return h.handleFirstStep(masterConn)
	})
	steps = append(steps, func() error {
		return h.handleSecondStep(masterConn)
	})
	steps = append(steps, func() error {
		return h.handleThirdStep(masterConn)
	})

	for _, fn := range steps {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (h *handshake) handleFirstStep(masterConn net.Conn) error {
	if err := h.sendPingCommand(masterConn); err != nil {
		return err
	}
	fmt.Println("PING")

	if err := h.verifyPingResponse(masterConn); err != nil {
		return err
	}

	fmt.Println("PONG")

	return nil
}

func (h *handshake) handleSecondStep(masterConn net.Conn) error {
	if err := h.sendListeningPortConfig(masterConn); err != nil {
		return err
	}

	fmt.Println("LISTENING PORT")

	if err := verifyOKResponse(masterConn); err != nil {
		return err
	}

	if err := h.sendCapabilityConfig(masterConn); err != nil {
		return err
	}

	fmt.Println("CAPA")

	if err := verifyOKResponse(masterConn); err != nil {
		return err
	}

	return nil
}

func (h *handshake) handleThirdStep(masterConn net.Conn) error {
	if err := h.sendPsyncCommand(masterConn); err != nil {
		return err
	}

	fmt.Println("PSYNC")

	if err := h.verifyPsyncResponse(masterConn); err != nil {
		return err
	}

	return nil
}

func (h *handshake) sendPsyncCommand(conn net.Conn) error {
	psyncCommand := []resp.SliceEncoding{
		{S: cmds.CMD_PSYNC, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "?", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "-1", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	_, err := conn.Write([]byte(resp.HandleEncodeSliceList(psyncCommand)))

	return err
}

func (h *handshake) verifyPingResponse(conn net.Conn) error {
	pingAnswer := []byte(resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_RESPONSE_PONG))
	response, err := getResponse(conn, len(pingAnswer))

	if err != nil || !bytes.Equal(pingAnswer, response) {
		return errors.New(cmds.CMD_RESPONSE_PONG + err.Error())
	}
	return nil
}

func (h *handshake) verifyPsyncResponse(conn net.Conn) error {
	_, err := getResponse(conn, 0)

	if err != nil {
		return errors.New(cmds.CMD_PSYNC + err.Error())
	}

	rdbFileBytes, err := getResponse(conn, 0)

	fmt.Print(string(rdbFileBytes))

	if err != nil {
		if collection, err := h.rdbReader.HandleReadFromBytes(rdbFileBytes); err != nil {
			h.inMemoryStorage = collection
		}
	}

	return err
}

func (h *handshake) sendCapabilityConfig(conn net.Conn) error {
	capabilityConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "capa", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "psync2", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}
	_, err := conn.Write([]byte(resp.HandleEncodeSliceList(capabilityConfig)))

	return err
}

func (h *handshake) sendPingCommand(conn net.Conn) error {
	_, err := conn.Write([]byte(resp.HandleEncodeSliceList([]resp.SliceEncoding{{S: cmds.CMD_PING, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING}})))

	return err
}

func (h *handshake) sendListeningPortConfig(conn net.Conn) error {
	listeningPortConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "listening-port", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: h.cfg.Port, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	_, err := conn.Write([]byte(resp.HandleEncodeSliceList(listeningPortConfig)))

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

func verifyOKResponse(conn net.Conn) error {
	okAnswer := []byte(resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_RESPONSE_OK))
	response, err := getResponse(conn, len(okAnswer))
	fmt.Println(string(response), "RESPONSE")

	if err != nil || !bytes.Equal(okAnswer, response) {
		return errors.New(cmds.CMD_RESPONSE_OK + err.Error())
	}
	return nil
}
