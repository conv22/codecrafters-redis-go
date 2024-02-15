package config

import (
	"flag"
)

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

func (config Config) IsReplica() bool {
	return config.Replication.Role == CONFIG_SLAVE_ROLE
}

func NewConfig() *Config {
	flag.Parse()
	flagArgs := flag.Args()
	return &Config{
		DirFlag:        *dirFlag,
		DbFilenameFlag: *dbFilenameFlag,
		Port:           *port,
		Replication:    NewReplicationInfo(*replica, flagArgs),
	}
}
