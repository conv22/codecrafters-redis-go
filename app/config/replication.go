package config

import "net"

type ReplicationInfo struct {
	Role          string
	Offset        string
	MasterReplId  string
	MasterAddress string
}

func NewReplicationInfo(replicaFlag string, flags []string) *ReplicationInfo {
	masterAddress := getMasterAddress(replicaFlag, flags)
	role := CONFIG_MASTER_ROLE
	if replicaFlag != "" {
		role = CONFIG_SLAVE_ROLE
	}
	return &ReplicationInfo{
		Role:          role,
		MasterAddress: masterAddress,
		MasterReplId:  "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Offset:        "0",
	}
}

func getMasterAddress(replicaFlag string, flags []string) (masterAddress string) {
	if len(flags) != 1 {
		return replicaFlag
	}

	return net.JoinHostPort(replicaFlag, flags[0])
}
