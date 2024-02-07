package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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

func parseResp(s string) ([]string, error) {
	clone := strings.Clone(s)
	const TERMINATOR = "\r\n"
	const ITEM_COUNT = "*"
	const ITEM_LENGTH = "$"

	first_item_index := strings.Index(clone, TERMINATOR)

	if first_item_index == -1 {
		return nil, errors.New("not valid string")
	}

	total_items_count, err := strconv.Atoi(clone[1:first_item_index])

	if err != nil {
		return nil, err
	}

	clone = clone[(first_item_index + len(TERMINATOR)):]
	curr_count := 0
	curr_item_length := 0
	words := make([]string, 0, total_items_count)

	for curr_count < total_items_count {
		if strings.HasPrefix(clone, ITEM_LENGTH) {
			first_item_index := strings.Index(clone, TERMINATOR)

			total_chars_count, err := strconv.Atoi(clone[1:first_item_index])
			if err != nil {
				return nil, err
			}
			clone = clone[first_item_index+len(TERMINATOR):]
			curr_item_length = total_chars_count
			continue
		}
		words = append(words, clone[0:curr_item_length])
		clone = clone[(len(TERMINATOR) + curr_item_length):]
		curr_count++
	}

	return words, nil

}

func handleClient(conn net.Conn, wg *sync.WaitGroup) {
	defer conn.Close()
	defer wg.Done()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	PING_RESPONSE := "+PONG\r\n"
	buffer := make([]byte, 1024)
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		input := string(buffer[:n])

		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "ping" {
			writer.WriteString(PING_RESPONSE)
		} else if values, err := parseResp(input); err == nil {
			if len(values) == 2 && strings.ToLower(values[0]) == "echo" {
				msg := "+" + values[1] + "\r\n"
				writer.WriteString(msg)
			}
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error flushing writer:", err)
			return
		}
	}
}
