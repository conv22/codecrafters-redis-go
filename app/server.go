package main

import (
	"bufio"
	"fmt"
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

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	parser := parser.CreateParser("resp")

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			break
		}

		line = strings.TrimSpace(line)

		write := func() string {
			if strings.EqualFold(line, "ping") {
				return "+PONG\r\n"
			}
			result, err := parser.HandleParse(line)
			// echo cmd
			if err == nil {
				if len(result) >= 2 {
					if strings.EqualFold(result[0], "echo") {
						return "+" + result[1] + "\r\n"
					}
				}
			}

			return line
		}()

		_, err = writer.WriteString(write)

		if err != nil {
			fmt.Println("Error writing to connection:", err)
			break
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing buffer:", err)
			break
		}
	}

}
