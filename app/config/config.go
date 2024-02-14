package config

import "flag"

type Config struct {
	DirFlag        string
	DbFilenameFlag string
	Port           string
}

var (
	dirFlag        = flag.String("dir", "", "The directory where RDB files are stored")
	dbFilenameFlag = flag.String("dbfilename", "", "The name of the RDB file")
	port           = flag.String("port", "6379", "Port for TCP server to listen to")
)

func InitializeConfig() *Config {

	return &Config{
		DirFlag:        *dirFlag,
		DbFilenameFlag: *dbFilenameFlag,
		Port:           *port,
	}
}
