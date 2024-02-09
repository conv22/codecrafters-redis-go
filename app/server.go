package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	parser "github.com/codecrafters-io/redis-starter-go/app/parsers"
)

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

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	parser := parser.CreateParser("resp")
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

		write := func() string {
			parserResult, err := parser.HandleParse(line)
			if err != nil || len(parserResult) == 0 {
				return line
			}

			firstCmd := parserResult[0]

			if strings.EqualFold(firstCmd, "ping") {
				return "+PONG\r\n"
			}

			if len(parserResult) >= 2 && strings.EqualFold(firstCmd, "echo") {
				return "+" + parserResult[1] + "\r\n"
			}

			return line
		}()
		writer.Write([]byte(write))
		err = writer.Flush()
		if err != nil {
			break
		}

	}
}
