package replication

import (
	"net"
	"sync"
)

type ReplicaClient struct {
	IsActive      bool
	ReplicationId string
	Offset        string
	connections   []net.Conn
	listeningPort string
	Mu            sync.Mutex
	// capa     string
	// psync2   string
}

func NewReplicaClient(listeningPort string) *ReplicaClient {
	return &ReplicaClient{
		connections:   []net.Conn{},
		listeningPort: listeningPort,
	}
}

func (client *ReplicaClient) AppendConnection(conn net.Conn) {
	client.Mu.Lock()
	defer client.Mu.Unlock()
	client.connections = append(client.connections, conn)
}
