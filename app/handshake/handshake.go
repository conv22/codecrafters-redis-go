package handshake

import (
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
	respReader      *resp.RespReader
	inMemoryStorage *storage.StorageCollection
	masterConn      net.Conn
	doneCh          chan struct{}
}

func New(cfg *config.Config, inMemoryStorage *storage.StorageCollection, masterConn net.Conn, respReader *resp.RespReader, doneCh chan struct{}) *handshake {
	return &handshake{
		cfg:             cfg,
		rdbReader:       rdb.NewRdb(),
		inMemoryStorage: inMemoryStorage,
		masterConn:      masterConn,
		respReader:      respReader,
		doneCh:          doneCh,
	}
}

func (h *handshake) HandleHandshake() error {
	steps := []func() error{h.handleFirstStep, h.handleSecondStep, h.handleThirdStep}

	for i, fn := range steps {
		if err := fn(); err != nil {
			return err
		}
		fmt.Printf("%d handshake step complete \n", i+1)
	}

	h.doneCh <- struct{}{}

	return nil
}

func (h *handshake) handleFirstStep() error {
	if err := h.sendPingCommand(); err != nil {
		return err
	}

	if err := h.verifyPingResponse(); err != nil {
		return err
	}

	return nil
}

func (h *handshake) handleSecondStep() error {
	if err := h.sendListeningPortConfig(); err != nil {
		return err
	}

	if err := h.verifyOKResponse(); err != nil {
		return err
	}

	if err := h.sendCapabilityConfig(); err != nil {
		return err
	}

	if err := h.verifyOKResponse(); err != nil {
		return err
	}

	return nil
}

func (h *handshake) handleThirdStep() error {
	if err := h.sendPsyncCommand(); err != nil {
		return err
	}

	if err := h.verifyPsyncResponse(); err != nil {
		return err
	}

	return nil
}

func (h *handshake) sendPsyncCommand() error {
	psyncCommand := []resp.SliceEncoding{
		{S: cmds.CMD_PSYNC, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "?", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "-1", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	_, err := h.masterConn.Write([]byte(resp.HandleEncodeSliceList(psyncCommand)))

	return err
}

func (h *handshake) verifyPingResponse() error {

	response, err := h.handleReadFromConnection()

	if err != nil {
		return err
	}

	if len(response) == 0 || cmds.CMD_RESPONSE_PONG != response[0].Value {
		return errors.New("invalid ping response")
	}

	return nil
}

func (h *handshake) verifyPsyncResponse() error {
	// ignore for now,  CMD_RESPONSE_FULL_RESYNC, h.replicationStore.MasterReplId, h.replicationStore.Offset
	_, err := h.handleReadFromConnection()

	if err != nil {
		return err
	}

	rdbFile, err := h.respReader.HandleReadRdbFile()

	if err != nil {
		return err
	}

	if collection, err := h.rdbReader.HandleReadFromBytes(rdbFile); err != nil {
		h.inMemoryStorage = collection
	}

	return err
}

func (h *handshake) sendCapabilityConfig() error {
	capabilityConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "capa", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "psync2", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}
	_, err := h.masterConn.Write([]byte(resp.HandleEncodeSliceList(capabilityConfig)))

	return err
}

func (h *handshake) sendPingCommand() error {
	_, err := h.masterConn.Write([]byte(resp.HandleEncodeSliceList([]resp.SliceEncoding{{S: cmds.CMD_PING, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING}})))

	return err
}

func (h *handshake) sendListeningPortConfig() error {
	listeningPortConfig := []resp.SliceEncoding{
		{S: cmds.CMD_REPLCONF, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: "listening-port", Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
		{S: h.cfg.Port, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING},
	}

	_, err := h.masterConn.Write([]byte(resp.HandleEncodeSliceList(listeningPortConfig)))

	return err
}

func (h *handshake) verifyOKResponse() error {
	response, err := h.handleReadFromConnection()

	if err != nil {
		return err
	}

	if len(response) == 0 || cmds.CMD_RESPONSE_OK != response[0].Value {
		return errors.New("OK validation failed")
	}

	return nil
}

func (h *handshake) handleReadFromConnection() ([]resp.ParsedCmd, error) {
	parsed, _, err := h.respReader.HandleRead()

	if err != nil {
		return nil, err
	}

	if len(parsed) > 0 {
		return parsed, nil
	}

	return nil, errors.New("no data received after multiple attempts")
}
