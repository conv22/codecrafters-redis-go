package replication

import (
	"flag"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	REPLICATION_MASTER_ROLE = "master"
	REPLICATION_SLAVE_ROLE  = "slave"
)

var replicaFlag = flag.String("replicaof", "", "The address for Master instance")

type replicaAddress = string

type ReplicationStore struct {
	Role            string
	offset          int64
	MasterReplId    string
	MasterAddress   string
	mu              sync.RWMutex
	AckCompleteChan chan replicaAddress
	ReplicasMap     map[replicaAddress]*ReplicaClient
}

func NewReplicationStore() *ReplicationStore {
	flag.Parse()

	role := determineRole()

	masterAddress := determineMasterAddress()

	return &ReplicationStore{
		Role:            role,
		MasterAddress:   masterAddress,
		MasterReplId:    "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		offset:          0,
		ReplicasMap:     make(map[replicaAddress]*ReplicaClient),
		AckCompleteChan: make(chan replicaAddress),
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
	return len(r.ReplicasMap) > 0
}

func (r *ReplicationStore) NumberOfReplicas() int {
	return len(r.ReplicasMap)
}

func (r *ReplicationStore) AppendClient(address string, client *ReplicaClient) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ReplicasMap[address] = client
}

func (r *ReplicationStore) PopulateCmdToReplicas(cmd []byte) {
	var wg sync.WaitGroup
	for _, replica := range r.ReplicasMap {
		wg.Add(1)
		go func(replica *ReplicaClient) {
			defer wg.Done()
			replica.PropagateCmd(cmd)

		}(replica)

	}
	wg.Wait()
}

func (r *ReplicationStore) GetAckFromReplicas() {
	cmd := []byte(resp.HandleEncodeSliceList([]resp.SliceEncoding{
		{
			S:        "REPLCONF",
			Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING,
		},
		{
			S:        "GETACK",
			Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING,
		},
		{
			S:        "*",
			Encoding: resp.RESP_ENCODING_CONSTANTS.BULK_STRING,
		},
	}))
	var wg sync.WaitGroup
	for _, replica := range r.ReplicasMap {
		if replica.offset == replica.expectedOffset {
			continue
		}
		wg.Add(1)
		go func(replica *ReplicaClient) {
			defer wg.Done()
			replica.PropagateCmd(cmd)
		}(replica)

	}
	wg.Wait()
}

func (r *ReplicationStore) GetNumOfAckReplicas() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	replicasInSync := 0

	for _, replica := range r.ReplicasMap {

		if replica.offset == replica.expectedOffset {
			replicasInSync++
		}
	}

	return replicasInSync
}

func GetReplicationAddress(conn net.Conn) (replicaAddress, error) {
	masterLocalAddr := conn.RemoteAddr().String()
	host, port, err := net.SplitHostPort(masterLocalAddr)
	if err != nil {
		return "", err
	}

	return net.JoinHostPort(host, port), nil
}

func (r *ReplicationStore) GetReplicaClientByAddress(address replicaAddress) (*ReplicaClient, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	replica, hasReplica := r.ReplicasMap[address]

	return replica, hasReplica
}

func (r *ReplicationStore) IncOffset(inc int64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.offset += inc
}

func (r *ReplicationStore) GetOffset() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.offset
}
