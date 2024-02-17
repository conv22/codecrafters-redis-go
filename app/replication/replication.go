package replication

import (
	"flag"
	"net"
)

type ReplicationInfo struct {
	Role          string
	Offset        string
	MasterReplId  string
	MasterAddress string
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

func (replication ReplicationInfo) IsReplica() bool {
	return replication.Role == REPLICATION_SLAVE_ROLE
}

func (replication ReplicationInfo) IsMaster() bool {
	return replication.Role == REPLICATION_MASTER_ROLE
}

func getMasterAddress(replicaFlag string, flags []string) (masterAddress string) {
	if len(flags) != 1 {
		return replicaFlag
	}

	return net.JoinHostPort(replicaFlag, flags[0])
}
