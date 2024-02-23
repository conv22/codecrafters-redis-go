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
var cmdProcessor = cmds.NewRespCmdProcessor(parser, inMemoryStorage, cfg, replicationStore)

var replicationStore = replication.NewReplicationStore()
var replicationChannel chan []byte

func initStorage() *storage.StorageCollection {
	persistStorage, err := rdbReader.HandleRead(path.Join(cfg.DirFlag, cfg.DbFilenameFlag))

	if err != nil {
		return storage.NewStorageCollection()
	}

	return persistStorage
}

func connectToMaster() (net.Conn, error) {
	masterConn, err := net.Dial("tcp", replicationStore.MasterAddress)
	if err != nil {
		return nil, errors.New("failed to connect to master: " + replicationStore.MasterAddress)
	}
	return masterConn, nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+cfg.Port)
	if err != nil {
		fmt.Println("Failed to bind to port: ", cfg.Port)
		os.Exit(1)
	}
	defer listener.Close()

	if replicationStore.IsReplica() {
		masterConn, err := connectToMaster()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		handleHandshake(masterConn)
		go handleConnection(masterConn, true)
	} else {
		replicationChannel = make(chan []byte)
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
		if !isPersistentConn && !replicationStore.IsReplicaClient(conn) {
			conn.Close()
		}
	}()

	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}

		bytesData := buf[:bytesRead]

		processedResult := cmdProcessor.ProcessCmd(bytesData, conn)

		for _, item := range processedResult {
			if len(item.Answer) > 0 {
				writer.Write([]byte(item.Answer))
			}

			if replicationChannel != nil && item.IsDuplicate {
				fmt.Printf("%s bytes input", item.BytesInput)
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
		replicationStore.PopulateCmdToReplicas(<-replicationChannel)
	}
}
