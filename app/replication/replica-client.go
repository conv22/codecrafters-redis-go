package replication

import (
	"bufio"
	"net"
)

type ReplicaClient struct {
	conn     net.Conn
	BuffCmds [][]byte
	Writer   *bufio.Writer
	// capa     string
	// psync2   string
}

func NewReplicaClient(conn net.Conn) *ReplicaClient {
	return &ReplicaClient{
		conn:   conn,
		Writer: bufio.NewWriter(conn),
	}
}
