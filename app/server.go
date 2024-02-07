package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!!!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

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
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	PING_RESPONSE := "+PONG\r\n"

	for scanner.Scan() {
		input := scanner.Text()

		if input == "PING" {
			writer.WriteString(PING_RESPONSE)
		} else if strings.HasPrefix(input, "ECHO") {
			result := strings.TrimPrefix(input, "ECHO")
			writer.WriteString(result + "\r\n")
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing writer:", err)
			return
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning user input")
	}

}
