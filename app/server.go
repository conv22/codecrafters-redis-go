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

	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
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
	RespEncodingConstants := parsers.RespEncodingConstants
	parser := parsers.CreateParser("resp")
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)
	tmpDb := make(map[string]string)

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
			parsedResult, err := parser.HandleParse(line)
			if err != nil || len(parsedResult) == 0 {
				return line
			}

			firstCmd := strings.ToLower(parsedResult[0])

			switch firstCmd {
			case "ping":
				return parser.HandleEncode(RespEncodingConstants.String, "PONG")

			case "echo":
				if len(parsedResult) >= 2 {
					return parser.HandleEncode(RespEncodingConstants.String, parsedResult[1])

				}

			case "set":
				if len(parsedResult) >= 3 {
					key, value := parsedResult[1], parsedResult[2]
					tmpDb[key] = value
					return parser.HandleEncode(RespEncodingConstants.String, "OK")
				}

			case "get":
				if len(parsedResult) >= 2 {
					key := parsedResult[1]
					value := tmpDb[key]
					return parser.HandleEncode(RespEncodingConstants.String, value)
				}

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
