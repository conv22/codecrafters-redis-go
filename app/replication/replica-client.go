package replication

import (
	"net"
	"sync"
)

type ReplicaClient struct {
	ReplicationId  string
	expectedOffset int64
	offset         int64
	conn           net.Conn
	listeningPort  string
	mu             sync.RWMutex
}

func (client *ReplicaClient) PropagateCmd(cmd []byte) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.expectedOffset = client.expectedOffset + int64(len(cmd))
	client.conn.Write(cmd)
}

func NewReplicaClient(listeningPort string, conn net.Conn) *ReplicaClient {
	return &ReplicaClient{
		conn:          conn,
		listeningPort: listeningPort,
	}
}

func (client *ReplicaClient) HandleAck(newOffset int64) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.offset = newOffset
}

func (client *ReplicaClient) InitailizeOffsetAndReplId(offset int64, replicationId string) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.offset = offset
	client.expectedOffset = offset
	client.ReplicationId = replicationId
}
