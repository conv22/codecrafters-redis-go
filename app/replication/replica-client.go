package replication

import (
	"net"
	"sync"
)

type ReplicaClient struct {
	ReplicationId string
	Offset        string
	connections   []net.Conn
	listeningPort string
	mu            sync.RWMutex
}

func NewReplicaClient(listeningPort string) *ReplicaClient {
	return &ReplicaClient{
		connections:   []net.Conn{},
		listeningPort: listeningPort,
	}
}

func (client *ReplicaClient) SetOffsetAndReplicationId(offset, replicationId string) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.Offset = offset
	client.ReplicationId = replicationId
}

func (client *ReplicaClient) AppendConnection(conn net.Conn) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.connections = append(client.connections, conn)
}
