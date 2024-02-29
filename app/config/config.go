package config

import (
	"flag"
	"sync"
)

type Config struct {
	DirFlag        string
	DbFilenameFlag string
	Port           string
	mu             sync.RWMutex
	offset         int64
}

var (
	dirFlag        = flag.String("dir", "", "The directory where RDB files are stored")
	dbFilenameFlag = flag.String("dbfilename", "", "The name of the RDB file")
	port           = flag.String("port", "6379", "Port for TCP server to listen to")
)

func NewConfig() *Config {
	flag.Parse()

	return &Config{
		DirFlag:        *dirFlag,
		DbFilenameFlag: *dbFilenameFlag,
		Port:           *port,
	}
}

func (cfg *Config) IncOffset(inc int64) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.offset += inc
}

func (cfg *Config) GetOffset() int64 {
	cfg.mu.RLock()
	defer cfg.mu.RUnlock()
	return cfg.offset
}
