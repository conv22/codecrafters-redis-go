package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path"

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
var cmdProcessor = cmds.NewRespCmdProcessor(parser, inMemoryStorage, cfg, replicationInfo)

var replicationInfo = replication.NewReplicationInfo()
var replicationChannel = make(chan []byte)

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
		masterConn, err := handleHandshake()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		go handleConnection(masterConn, true)
	} else {
		go handleSyncWithReplicas()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn, false)
	}
}

func handleConnection(conn net.Conn, isPersistentConn bool) {
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)

	defer func() {
		if !isPersistentConn && !replicationInfo.IsReplicaClient(conn) {
			conn.Close()
		}
	}()

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

		processedResult := cmdProcessor.ProcessCmd(bytesData, conn)

		for _, item := range processedResult {
			if len(item.Answer) > 0 {
				writer.Write([]byte(item.Answer))
			}

			if item.IsDuplicate && replicationInfo.IsMaster() {
				replicationChannel <- item.BytesInput
			}
		}

		if err := writer.Flush(); err != nil {
			break
		}

	}

}

func handleSyncWithReplicas() {
	for {
		cmds := <-replicationChannel
		replicationInfo.PopulateCmdToReplicas(cmds)
	}
}
