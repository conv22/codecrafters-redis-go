package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	cmds "github.com/codecrafters-io/redis-starter-go/app/cmds"
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

var tmpDb = storage.CreateStorage()

func main() {
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
		go handleClient(conn, &wg)
	}
}

type db_record struct {
	value          string
	expirationTime *int64
}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()
	parser := parsers.CreateParser("resp")
	cmdProcessor := cmds.CreateProcessor("resp", &parser, &tmpDb)
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

		result, err := cmdProcessor.ProcessCmd(line)

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
