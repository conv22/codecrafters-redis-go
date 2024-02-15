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
var cmdProcessor = cmds.NewRespCmdProcessor(parser, inMemoryStorage, cfg)

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
		err := handleHandshake()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		wg.Add(1)
		go handleClient(conn, &wg)
	}
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
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
		if err := writer.Flush(); err != nil {
			break
		}

	}
}
