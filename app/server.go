package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/rdb"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

var cfg = config.NewConfig()
var rdbReader = rdb.NewRdb()
var inMemoryStorage = initStorage()
var parser = resp.NewRespParser()
var replicationInfo = replication.NewReplicationInfo()

var replicationLock sync.Mutex

func initStorage() *storage.StorageCollection {
	persistStorage, err := rdbReader.HandleRead(path.Join(cfg.DirFlag, cfg.DbFilenameFlag))

	if err != nil {
		return storage.NewStorageCollection()
	}

	return persistStorage
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+cfg.Port)
	if err != nil {
		fmt.Println("Failed to bind to port: ", cfg.Port)
		os.Exit(1)
	}
	defer listener.Close()

	if replicationInfo.IsReplica() {
		err := handleHandshake()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)
	cmdProcessor := cmds.NewRespCmdProcessor(parser, inMemoryStorage, cfg, replicationInfo, conn)

	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading from connection:", err)
			break
		}

		bytesData := buf[:bytesRead]
		line := string(bytesData)

		processedResult := cmdProcessor.ProcessCmd(line)

		for _, item := range processedResult {
			if len(item.Answer) > 0 {
				writer.Write([]byte(item.Answer))
			}

			if item.IsDuplicate && replicationInfo.IsMaster() {
				go handleSyncWithReplica(bytesData)
			}
		}

		if err := writer.Flush(); err != nil {
			break
		}

	}

	if !replicationInfo.IsReplicaClient(conn) {
		conn.Close()
	}

}

func handleSyncWithReplica(cmd []byte) {
	replicationLock.Lock()

	// TODO: send cmds only once replica is active
	for _, replica := range replicationInfo.Replicas {
		replica.Writer.Write(cmd)
		replica.Writer.Flush()
	}
	replicationLock.Unlock()

}
