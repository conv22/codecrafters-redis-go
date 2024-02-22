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
	Mu            sync.Mutex
	Replicas      map[string]*ReplicaClient
}

var replica = flag.String("replicaof", "", "The address for Master instance")

func NewReplicationInfo() *ReplicationInfo {
	flag.Parse()
	// Todo: handle flags differently
	flagArgs := flag.Args()
	masterAddress := getMasterAddress(*replica, flagArgs)
	role := REPLICATION_MASTER_ROLE
	if *replica != "" {
		role = REPLICATION_SLAVE_ROLE
	}
	return &ReplicationInfo{
		Role:          role,
		MasterAddress: masterAddress,
		MasterReplId:  "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Offset:        "0",
		Replicas:      make(map[string]*ReplicaClient),
	}
}

func (replication *ReplicationInfo) IsReplica() bool {
	return replication.Role == REPLICATION_SLAVE_ROLE
}

func (replication *ReplicationInfo) IsMaster() bool {
	return replication.Role == REPLICATION_MASTER_ROLE
}

func (replication *ReplicationInfo) AppendClient(address string, client *ReplicaClient) {
	replication.Mu.Lock()
	replication.Replicas[address] = client
	replication.Mu.Unlock()
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
	replication.Mu.Lock()
	defer replication.Mu.Unlock()
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

func (replication *ReplicationInfo) IsReplicaMaster(conn net.Conn) bool {
	if replication.IsMaster() {
		return false
	}
	connAddress, err := GetReplicationAddress(conn)

	if err != nil {
		return false
	}

	return replication.MasterAddress == connAddress

}

func getMasterAddress(replicaFlag string, flags []string) (masterAddress string) {
	if len(flags) != 1 {
		return replicaFlag
	}

	return net.JoinHostPort(replicaFlag, flags[0])
}
