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
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

var cfg = config.NewConfig()
var rdbReader = rdb.NewRdb()
var inMemoryStorage = initStorage()
var parser = resp.NewRespParser()

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

	if cfg.IsReplica() {
		replication := cfg.Replication
		masterConn, err := net.Dial("tcp", replication.MasterAddress)
		if err != nil {
			fmt.Println("Master is not available", replication.MasterAddress)
			os.Exit(1)
		}
		defer masterConn.Close()
		masterConn.Write([]byte(parser.HandleEncodeSliceList([]resp.SliceEncoding{
			{
				S:        "ping",
				Encoding: resp.RESP_ENCODING_CONSTANTS.STRING,
			},
		})))
	}

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		wg.Add(1)
		go handleClient(conn, &wg, cfg)
	}
}

func handleClient(conn net.Conn, wg *sync.WaitGroup, config *config.Config) {
	defer conn.Close()
	defer wg.Done()
	cmdProcessor := cmds.NewRespCmdProcessor(parser, inMemoryStorage, config)
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)

	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading from connection:", err)
			break
		}

		line := string(buf[:bytesRead])

		result := cmdProcessor.ProcessCmd(line)

		write := line

		if err == nil {
			write = result
		}
		writer.Write([]byte(write))
		err = writer.Flush()
		if err != nil {
			break
		}
	}
}
