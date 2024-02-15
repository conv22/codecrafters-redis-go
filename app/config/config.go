package config

import "flag"

type ReplicationInfo struct {
	Role         string
	Offset       string
	MasterReplId string
}

type Config struct {
	DirFlag        string
	DbFilenameFlag string
	Port           string
	Replication    *ReplicationInfo
}

var (
	dirFlag        = flag.String("dir", "", "The directory where RDB files are stored")
	dbFilenameFlag = flag.String("dbfilename", "", "The name of the RDB file")
	port           = flag.String("port", "6379", "Port for TCP server to listen to")
	replica        = flag.String("replicaof", "", "The address for Master instance")
)

func getRole(masterAddress string) string {
	if masterAddress == "" {
		return "master"
	}
	return "slave"
}

func NewConfig() *Config {
	flag.Parse()
	return &Config{
		DirFlag:        *dirFlag,
		DbFilenameFlag: *dbFilenameFlag,
		Port:           *port,
		Replication: &ReplicationInfo{
			Role:         getRole(*replica),
			MasterReplId: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
			Offset:       "0",
		},
	}
}
