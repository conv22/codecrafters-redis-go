package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	cmds "github.com/codecrafters-io/redis-starter-go/app/cmds"
	config "github.com/codecrafters-io/redis-starter-go/app/config"
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

var inMemoryDb = storage.NewInMemoryStorage()

func main() {
	flag.Parse()

	cfg := config.InitializeConfig()

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

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
	parser := parsers.NewRespParser()
	cmdProcessor := cmds.NewRespCmdProcessor(parser, inMemoryDb, config)
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
