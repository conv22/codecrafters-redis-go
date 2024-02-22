package replication

import (
	"bufio"
	"net"
	"sync"
)

type ReplicaClient struct {
	IsActive      bool
	ReplicationId string
	Offset        string
	conn          net.Conn
	Writer        *bufio.Writer
	listeningPort string
	Mu            sync.Mutex
	// capa     string
	// psync2   string
}

func NewReplicaClient(conn net.Conn, listeningPort string) *ReplicaClient {
	return &ReplicaClient{
		conn:          conn,
		Writer:        bufio.NewWriter(conn),
		listeningPort: listeningPort,
	}
}
