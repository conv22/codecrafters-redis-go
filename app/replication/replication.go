package replication

import (
	"flag"
	"net"
	"sync"
)

type ReplicationInfo struct {
	Role          string
	Offset        string
	MasterReplId  string
	MasterAddress string
	mu            sync.Mutex
	Replicas      map[string]*ReplicaClient
}

var replicaFlag = flag.String("replicaof", "", "The address for Master instance")

func NewReplicationInfo() *ReplicationInfo {
	flag.Parse()

	role := determineRole()

	masterAddress := determineMasterAddress()

	return &ReplicationInfo{
		Role:          role,
		MasterAddress: masterAddress,
		MasterReplId:  "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Offset:        "0",
		Replicas:      make(map[string]*ReplicaClient),
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

func (replication *ReplicationInfo) IsReplica() bool {
	return replication.Role == REPLICATION_SLAVE_ROLE
}

func (replication *ReplicationInfo) IsMaster() bool {
	return replication.Role == REPLICATION_MASTER_ROLE
}

func (replication *ReplicationInfo) AppendClient(address string, client *ReplicaClient) {
	replication.mu.Lock()
	replication.Replicas[address] = client
	replication.mu.Unlock()
}

func GetReplicationAddress(conn net.Conn) (string, error) {
	masterLocalAddr := conn.LocalAddr().String()
	host, port, err := net.SplitHostPort(masterLocalAddr)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(host, port), nil
}

func (replication *ReplicationInfo) PopulateCmdToReplicas(data []byte) {
	replication.mu.Lock()
	defer replication.mu.Unlock()
	for _, replica := range replication.Replicas {
		for _, conn := range replica.connections {
			conn.Write(data)
		}
	}
}

func (replication *ReplicationInfo) IsReplicaClient(conn net.Conn) bool {
	connAddress, err := GetReplicationAddress(conn)

	if err != nil {
		return false
	}

	_, hasReplica := replication.Replicas[connAddress]

	return hasReplica
}
