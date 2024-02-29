package replication

import (
	"flag"
	"net"
	"sync"
)

const (
	REPLICATION_MASTER_ROLE = "master"
	REPLICATION_SLAVE_ROLE  = "slave"
)

var replicaFlag = flag.String("replicaof", "", "The address for Master instance")

type ReplicationStore struct {
	Role          string
	Offset        string
	MasterReplId  string
	MasterAddress string
	mu            sync.RWMutex
	queueMu       sync.RWMutex
	replicasMap   map[string]*ReplicaClient
}

func NewReplicationStore() *ReplicationStore {
	flag.Parse()

	role := determineRole()

	masterAddress := determineMasterAddress()

	return &ReplicationStore{
		Role:          role,
		MasterAddress: masterAddress,
		MasterReplId:  "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Offset:        "0",
		replicasMap:   make(map[string]*ReplicaClient),
	}
}

func determineRole() string {
	if *replicaFlag != "" {
		return REPLICATION_SLAVE_ROLE
	}
	return REPLICATION_MASTER_ROLE
}

func determineMasterAddress() string {
	flagArgs := flag.Args()
	if len(flagArgs) == 1 {
		return net.JoinHostPort(*replicaFlag, flagArgs[0])
	}
	return *replicaFlag
}

func (r *ReplicationStore) IsReplica() bool {
	return r.Role == REPLICATION_SLAVE_ROLE
}

func (r *ReplicationStore) IsMaster() bool {
	return r.Role == REPLICATION_MASTER_ROLE
}

func (r *ReplicationStore) HasReplicas() bool {
	return len(r.replicasMap) > 0
}

func (r *ReplicationStore) NumberOfReplicas() int {
	return len(r.replicasMap)
}

func (r *ReplicationStore) AppendClient(address string, client *ReplicaClient) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.replicasMap[address] = client
}

func (r *ReplicationStore) PopulateCmdToReplicas(cmd []byte) {
	r.queueMu.Lock()
	defer r.queueMu.Unlock()
	var wg sync.WaitGroup
	for _, replica := range r.replicasMap {
		replica.mu.Lock()
		wg.Add(len(replica.connections))
		for _, conn := range replica.connections {
			go func(conn net.Conn) {
				defer wg.Done()
				conn.Write(cmd)
			}(conn)
		}
		replica.mu.Unlock()
	}
	wg.Wait()
}
func GetReplicationAddress(conn net.Conn) (string, error) {
	masterLocalAddr := conn.LocalAddr().String()
	host, port, err := net.SplitHostPort(masterLocalAddr)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(host, port), nil
}

func (r *ReplicationStore) GetReplicaClientByAddress(address string) (*ReplicaClient, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	replica, hasReplica := r.replicasMap[address]

	return replica, hasReplica
}

func (r *ReplicationStore) IsReplicaClient(conn net.Conn) bool {
	connAddress, err := GetReplicationAddress(conn)

	if err != nil {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, hasReplica := r.replicasMap[connAddress]

	return hasReplica
}
