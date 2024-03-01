package replication

import (
	"net"
	"sync"
)

type ReplicaClient struct {
	ReplicationId  string
	Offset         int64
	conn           net.Conn
	expectedOffset int64
	isOffsetInSync bool
	listeningPort  string
	mu             sync.RWMutex
}

func (client *ReplicaClient) PropagateCmd(cmd []byte) {
	defer func() {
		client.expectedOffset += int64(len(cmd))
		client.isOffsetInSync = false
	}()
	client.conn.Write(cmd)
}

func NewReplicaClient(listeningPort string, conn net.Conn) *ReplicaClient {
	return &ReplicaClient{
		conn:          conn,
		listeningPort: listeningPort,
	}
}

func (client *ReplicaClient) SetOffsetAndReplicationId(offset int64, replicationId string) {
	client.mu.Lock()
	defer client.mu.Unlock()
	client.Offset = offset
	client.ReplicationId = replicationId
}
