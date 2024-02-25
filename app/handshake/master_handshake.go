package handshake

import (
	"bytes"
	"errors"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type masterHandshake struct {
	cfg             *config.Config
	rdbReader       *rdb.Rdb
	inMemoryStorage *storage.StorageCollection
}

func NewMasterHandshake(cfg *config.Config, inMemoryStorage *storage.StorageCollection) *masterHandshake {
	return &masterHandshake{
		cfg:             cfg,
		rdbReader:       rdb.NewRdb(),
		inMemoryStorage: inMemoryStorage,
	}
}

func (h *masterHandshake) HandleHandshake(masterConn net.Conn) error {
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

func (h *masterHandshake) handleFirstStep(masterConn net.Conn) error {
	if err := h.sendPingCommand(masterConn); err != nil {
		return err
	}

	if err := h.verifyPingResponse(masterConn); err != nil {
		return err
	}

	return nil
}

func (h *masterHandshake) handleSecondStep(masterConn net.Conn) error {
	if err := h.sendListeningPortConfig(masterConn); err != nil {
		return err
	}

	if err := verifyOKResponse(masterConn); err != nil {
		return err
	}

	if err := h.sendCapabilityConfig(masterConn); err != nil {
		return err
	}

	if err := verifyOKResponse(masterConn); err != nil {
		return err
	}

	return nil
}

func (h *masterHandshake) handleThirdStep(masterConn net.Conn) error {
	if err := h.sendPsyncCommand(masterConn); err != nil {
		return err
	}

	if err := h.verifyPsyncResponse(masterConn); err != nil {
		return err
	}

	return nil
}

func (h *masterHandshake) sendPsyncCommand(conn net.Conn) error {
	psyncCommand := []resp.SliceEncoding{
		{S: HANDSHAKE_CMD_PSYNC, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "?", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "-1", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	if err := sendSliceCommand(conn, psyncCommand); err != nil {
		return err
	}

	return nil
}

func (h *masterHandshake) verifyPingResponse(conn net.Conn) error {
	pingAnswer := []byte(resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, cmds.CMD_RESPONSE_PONG))
	response, err := getResponse(conn, len(pingAnswer))

	if err != nil || !bytes.Equal(pingAnswer, response) {
		return errors.New(cmds.CMD_RESPONSE_PONG + err.Error())
	}
	return nil
}

func (h *masterHandshake) verifyPsyncResponse(conn net.Conn) error {
	_, err := getResponse(conn, 1024)

	if err != nil {
		return errors.New(HANDSHAKE_CMD_PSYNC + err.Error())
	}

	rdbFileBytes, err := getResponse(conn, 1024)

	if err != nil {
		return err
	}

	if err != nil {
		collection, err := h.rdbReader.HandleReadFromBytes(rdbFileBytes)
		if err != nil {
			h.inMemoryStorage = collection
		}
	}

	return err
}

func (h *masterHandshake) sendCapabilityConfig(conn net.Conn) error {
	capabilityConfig := []resp.SliceEncoding{
		{S: HANDSHAKE_CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "capa", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "psync2", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	return sendSliceCommand(conn, capabilityConfig)
}

func (h *masterHandshake) sendPingCommand(conn net.Conn) error {
	return sendSliceCommand(conn, []resp.SliceEncoding{{S: cmds.CMD_PING, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING}})
}

func (h *masterHandshake) sendListeningPortConfig(conn net.Conn) error {
	listeningPortConfig := []resp.SliceEncoding{
		{S: HANDSHAKE_CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "listening-port", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: h.cfg.Port, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	return sendSliceCommand(conn, listeningPortConfig)
}
