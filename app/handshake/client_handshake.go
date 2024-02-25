package handshake

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type clientHandshake struct {
	cfg              *config.Config
	rdbReader        *rdb.Rdb
	inMemoryStorage  *storage.StorageCollection
	replicationStore *replication.ReplicationStore
}

func NewClientHandshake(cfg *config.Config, inMemoryStorage *storage.StorageCollection, replicationStore *replication.ReplicationStore) *clientHandshake {
	return &clientHandshake{
		cfg:              cfg,
		rdbReader:        rdb.NewRdb(),
		inMemoryStorage:  inMemoryStorage,
		replicationStore: replicationStore,
	}
}

// Starts from repl conf cmd
func (h *clientHandshake) HandleHandshake(clientConn net.Conn, parsed []resp.ParsedCmd) error {
	var steps []func() error

	steps = append(steps, func() error {
		return h.handleFirstStep(parsed, clientConn)
	})
	steps = append(steps, func() error {
		return h.handleSecondStep(clientConn)
	})

	for _, fn := range steps {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
}

func (h *clientHandshake) handleFirstStep(parsed []resp.ParsedCmd, clientConn net.Conn) error {
	replicationAddress, err := replication.GetReplicationAddress(clientConn)
	if err != nil {
		return errors.New("invalid connection address")
	}

	listeningPort := parsed[1].Value

	client, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		client = replication.NewReplicaClient(listeningPort)
		h.replicationStore.AppendClient(replicationAddress, client)
	}
	client.AppendConnection(clientConn)

	if err := sendOkCommand(clientConn); err != nil {
		return err
	}

	_, err = getResponse(clientConn, 1024)

	if err != nil {
		return err
	}

	if err := sendOkCommand(clientConn); err != nil {
		return err
	}

	return nil
}

func (h *clientHandshake) handleSecondStep(clientConn net.Conn) error {

	response, err := getResponse(clientConn, 1024)

	if err != nil {
		return err
	}

	parsed, err := resp.HandleParse(string(response))

	if err != nil || len(parsed) < 1 || len(parsed[0]) < 1 {
		return errors.New("invalid command")
	}

	replicationAddress, err := replication.GetReplicationAddress(clientConn)
	if err != nil {
		return errors.New("invalid connection address")
	}

	replica, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return errors.New("invalid connection address")
	}

	replica.SetOffsetAndReplicationId(parsed[0][0].Value, parsed[0][1].Value)

	const EMPTY_DB_HEX string = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

	decoded, err := hex.DecodeString(EMPTY_DB_HEX)
	if err != nil {
		return nil
	}

	ackCmd := resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.STRING, fmt.Sprintf("%s %s %s", HANDSHAKE_CMD_RESPONSE_FULL_RESYNC, h.replicationStore.MasterReplId, h.replicationStore.Offset))

	if err := h.sendAcknowledgmentCommand(clientConn, ackCmd); err != nil {
		return err
	}

	if err := verifyOKResponse(clientConn); err != nil {
		return err
	}

	encodingCmd := resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.BULK_STRING, string(decoded))

	if err := h.sendEncodingCommand(clientConn, encodingCmd); err != nil {
		return err
	}

	if err := verifyOKResponse(clientConn); err != nil {
		return err
	}

	return nil
}

func (h *clientHandshake) sendAcknowledgmentCommand(conn net.Conn, ackCmd string) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	ackSlice := []resp.SliceEncoding{{S: ackCmd, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING}}

	return sendCommand(writer, ackSlice)
}

func (h *clientHandshake) sendEncodingCommand(conn net.Conn, encodingCmd string) error {
	writer := bufio.NewWriter(conn)
	defer writer.Flush()

	encodingSlice := []resp.SliceEncoding{{S: encodingCmd, Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING}}

	return sendCommand(writer, encodingSlice)
}
