package config

import "flag"

type Config struct {
	DirFlag        string
	DbFilenameFlag string
	Port           string
	Role           string
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
		Role:           getRole(*replica),
	}
}
