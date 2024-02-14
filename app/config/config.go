package config

import "flag"

type Config struct {
	DirFlag        string
	DbFilenameFlag string
}

var (
	dirFlag        = flag.String("dir", "", "The directory where RDB files are stored")
	dbFilenameFlag = flag.String("dbfilename", "", "The name of the RDB file")
)

func InitializeConfig() *Config {
	return &Config{
		DirFlag:        *dirFlag,
		DbFilenameFlag: *dbFilenameFlag,
	}
}
