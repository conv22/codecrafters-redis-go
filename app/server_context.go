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
	inMemoryStorage  *storage.StorageCollection
	replicationStore *replication.ReplicationStore
}

func NewServerContext() *ServerContext {
	cfg := config.NewConfig()
	inMemoryStorage := initStorage(rdb.NewRdb(), cfg)
	replicationStore := replication.NewReplicationStore()
	return &ServerContext{
		cfg:              cfg,
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
