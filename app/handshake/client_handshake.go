package handshake

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type clientHandshake struct {
	cfg              *config.Config
	inMemoryStorage  *storage.StorageCollection
	replicationStore *replication.ReplicationStore
}

func NewClientHandshake(cfg *config.Config, inMemoryStorage *storage.StorageCollection, replicationStore *replication.ReplicationStore) *clientHandshake {
	return &clientHandshake{
		cfg:              cfg,
		inMemoryStorage:  inMemoryStorage,
		replicationStore: replicationStore,
	}
}

func (h *clientHandshake) HandleHandshake(clientConn net.Conn, parsedCmd []resp.ParsedCmd) (nextCmds []resp.ParsedCmd, err error) {
	handshakeErr := h.handleHandshakeStep(parsedCmd, clientConn)
	for handshakeErr == nil {
		nextCmd, err := getResponse(clientConn, 1024)
		if err != nil {
			return nil, err
		}
		parsedCmd, err = resp.HandleParse(string(nextCmd))
		if err != nil {
			return nil, err
		}
		nextCmds = parsedCmd
		handshakeErr = h.handleHandshakeStep(parsedCmd, clientConn)
	}

	if errors.Is(handshakeErr, errUnknownCmd) {
		return nextCmds, nil
	}

	return nil, handshakeErr
}

var errUnknownCmd = errors.New("unknown handshake command")

func (h *clientHandshake) handleHandshakeStep(parsedCmd []resp.ParsedCmd, clientConn net.Conn) error {
	switch strings.ToUpper(parsedCmd[0].Value) {
	case HANDSHAKE_CMD_REPLCONF:
		if err := h.handleReplConfCommand(parsedCmd[1:], clientConn); err != nil {
			return err
		}
	case HANDSHAKE_CMD_PSYNC:
		if err := h.handlePsyncCommand(parsedCmd[1:], clientConn); err != nil {
			return err
		}
	default:
		return errUnknownCmd
	}

	return nil
}

func (h *clientHandshake) handleReplConfCommand(args []resp.ParsedCmd, clientConn net.Conn) error {
	switch strings.ToUpper(args[0].Value) {
	case HANDSHAKE_CMD_LISTENING_PORT:
		if len(args) >= 2 {
			return h.handleCreateClient(args[1].Value, clientConn)
		}
	case HANDSHAKE_CMD_CAPA:
		if len(args) >= 2 {
			return h.handleReplCapa(args[1].Value, clientConn)
		}
	}
	return errors.New("invalid REPLCONF command: missing value")
}

func (h *clientHandshake) handleCreateClient(listeningPort string, clientConn net.Conn) error {
	replicationAddress, err := replication.GetReplicationAddress(clientConn)
	if err != nil {
		return errors.New("invalid connection address")
	}

	client, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		client = replication.NewReplicaClient(listeningPort)
		h.replicationStore.AppendClient(replicationAddress, client)
	}
	client.AppendConnection(clientConn)

	if err := sendOkCommand(clientConn); err != nil {
		return err
	}

	return nil
}

func (h *clientHandshake) handleReplCapa(value string, clientConn net.Conn) error {
	if err := sendOkCommand(clientConn); err != nil {
		return err
	}

	return nil
}

func (h *clientHandshake) handlePsyncCommand(args []resp.ParsedCmd, clientConn net.Conn) error {
	replicationAddress, err := replication.GetReplicationAddress(clientConn)
	if err != nil {
		return errors.New("invalid connection address")
	}

	replica, ok := h.replicationStore.GetReplicaClientByAddress(replicationAddress)
	if !ok {
		return errors.New("invalid connection address")
	}

	replica.SetOffsetAndReplicationId(args[0].Value, args[1].Value)

	const EMPTY_DB_HEX string = "524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2"

	decoded, err := hex.DecodeString(EMPTY_DB_HEX)
	if err != nil {
		return err
	}

	if err := h.sendAcknowledgmentCommand(clientConn, fmt.Sprintf("%s %s %s", HANDSHAKE_CMD_RESPONSE_FULL_RESYNC, h.replicationStore.MasterReplId, h.replicationStore.Offset)); err != nil {
		return err
	}

	if err := h.sendEncodingCommand(clientConn, string(decoded)); err != nil {
		return err
	}

	if err := verifyOKResponse(clientConn); err != nil {
		return err
	}

	return nil
}

func (h *clientHandshake) sendAcknowledgmentCommand(conn net.Conn, ackCmd string) error {
	return sendCommand(conn, resp.SliceEncoding{S: ackCmd, Encoding: resp.RESP_ENCODING_CONSTANTS.STRING})
}

func (h *clientHandshake) sendEncodingCommand(conn net.Conn, encodingCmd string) error {
	encodingCmd = resp.HandleEncode(resp.RESP_ENCODING_CONSTANTS.BULK_STRING, encodingCmd)
	encodingCmd = strings.TrimSuffix(encodingCmd, resp.RESP_ENCODING_CONSTANTS.SEPARATOR)
	_, err := conn.Write([]byte(encodingCmd))
	return err
}
