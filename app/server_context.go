package main

import (
	"path"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

type ServerContext struct {
	cfg              *config.Config
	rdbReader        *rdb.Rdb
	inMemoryStorage  *storage.StorageCollection
	replicationStore *replication.ReplicationStore
}

func NewServerContext() *ServerContext {
	cfg := config.NewConfig()
	rdbReader := rdb.NewRdb()
	inMemoryStorage := initStorage(rdbReader, cfg)
	replicationStore := replication.NewReplicationStore()
	return &ServerContext{
		cfg:              cfg,
		rdbReader:        rdbReader,
		inMemoryStorage:  inMemoryStorage,
		replicationStore: replicationStore,
	}

}

func initStorage(r *rdb.Rdb, c *config.Config) *storage.StorageCollection {
	persistStorage, err := r.HandleReadFromFile(path.Join(c.DirFlag, c.DbFilenameFlag))

	if err != nil {
		return storage.NewStorageCollection()
	}

	return persistStorage
}

var serverContext = NewServerContext()
